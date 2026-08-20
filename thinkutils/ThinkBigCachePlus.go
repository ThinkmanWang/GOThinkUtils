package thinkutils

import (
	"bufio"
	"encoding/gob"
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron"
)

// thinkBigCachePlusDiskEntry 是磁盘持久化的单条记录。
// 逐条流式 gob 编码，任意时刻内存里只有一条记录，避免一次性把整份数据拷进 buffer 造成的瞬时分配与 GC 压力。
// Partition 字段用于落盘/加载时按分区归位数据。
type thinkBigCachePlusDiskEntry struct {
	Partition string
	Key       string
	Value     []byte
}

type OnThinkBigCachePlusUpdate func()

// ThinkBigCachePlusToByte 把内存中的实际 struct 编码为 []byte 用于落盘。
type ThinkBigCachePlusToByte func(pData any) ([]byte, error)

// ThinkBigCachePlusFromByte 把磁盘上的 []byte 还原为实际 struct 存入内存。
// 与 ThinkBigCachePlusToByte 成对注册，保证 loadFromDisk 后 Get 直接拿到 struct，避免每次读取解码带来的 GC 压力。
type ThinkBigCachePlusFromByte func(data []byte) (any, error)

type ThinkBigCachePlusUpdateType int

const (
	THINK_BIGCACHE_PLUS_UPDATE_1_MIN ThinkBigCachePlusUpdateType = iota
	THINK_BIGCACHE_PLUS_UPDATE_5_MIN
	THINK_BIGCACHE_PLUS_UPDATE_10_MIN
	THINK_BIGCACHE_PLUS_UPDATE_30_MIN
	THINK_BIGCACHE_PLUS_UPDATE_1_HOUR
	THINK_BIGCACHE_PLUS_UPDATE_6_HOUR
	THINK_BIGCACHE_PLUS_UPDATE_12_HOUR
	THINK_BIGCACHE_PLUS_UPDATE_1_DAY
)

type ThinkBigCachePlusPartition struct {
	m_mapData      sync.Map
	m_funcToByte   ThinkBigCachePlusToByte
	m_funcFromByte ThinkBigCachePlusFromByte
}

type ThinkBigCachePlus struct {
	m_bStarted     atomic.Bool
	m_bSavedToDisk atomic.Bool
	m_lock         sync.RWMutex
	m_lockFile     sync.Mutex
	m_szFileName   string

	m_mapPartition    sync.Map
	m_nSaveToDiskType ThinkBigCachePlusUpdateType

	m_pCronJobs *gocron.Scheduler

	m_lst1MinListener   []OnThinkBigCachePlusUpdate
	m_lst5MinListener   []OnThinkBigCachePlusUpdate
	m_lst10MinListener  []OnThinkBigCachePlusUpdate
	m_lst30MinListener  []OnThinkBigCachePlusUpdate
	m_lst1HourListener  []OnThinkBigCachePlusUpdate
	m_lst6HourListener  []OnThinkBigCachePlusUpdate
	m_lst12HourListener []OnThinkBigCachePlusUpdate
	m_lst1DayListener   []OnThinkBigCachePlusUpdate
}

var (
	g_pThinkBigCachePlusInstance *ThinkBigCachePlus = &ThinkBigCachePlus{
		m_szFileName: "ThinkBigCachePlus.data",
	}
)

func ThinkBigCachePlusInstance() *ThinkBigCachePlus {
	return g_pThinkBigCachePlusInstance
}

// getOrCreatePartition 返回指定分区(不存在则原子创建)。
// 存指针而非值：ThinkBigCachePlusPartition 内含 sync.Map，绝不能被拷贝。
func (this *ThinkBigCachePlus) getOrCreatePartition(szPartition string) *ThinkBigCachePlusPartition {
	if v, ok := this.m_mapPartition.Load(szPartition); ok {
		return v.(*ThinkBigCachePlusPartition)
	}
	actual, _ := this.m_mapPartition.LoadOrStore(szPartition, &ThinkBigCachePlusPartition{})
	return actual.(*ThinkBigCachePlusPartition)
}

// loadFromDisk 从本地文件逐条流式解码，按分区把数据还原成实际 struct 存入 m_mapPartition。
// 注意：需在 Start 之前先通过 RegToByteFunction 注册各分区的编解码函数；
// 若某分区在加载时尚未注册 FromByte，则跳过该分区的数据(记录告警)。
func (this *ThinkBigCachePlus) loadFromDisk() error {
	this.m_lockFile.Lock()
	defer this.m_lockFile.Unlock()

	if _, err := os.Stat(this.m_szFileName); os.IsNotExist(err) {
		return nil
	}

	f, err := os.Open(this.m_szFileName)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer f.Close()

	// 逐条流式解码：任意时刻内存里只有一条记录，避免一次性载入整份数据。
	dec := gob.NewDecoder(bufio.NewReaderSize(f, 1<<20))
	for {
		var entry thinkBigCachePlusDiskEntry
		if err := dec.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			// 兼容旧格式或文件损坏：记录后按当前已加载数据处理，绝不阻断服务启动。
			log.Error("ThinkBigCachePlus loadFromDisk decode failed: %s", err.Error())
			return nil
		}

		p := this.getOrCreatePartition(entry.Partition)
		fromByte := p.m_funcFromByte
		if nil == fromByte {
			// 该分区未注册 FromByte，无法还原 struct，跳过。
			log.Error("ThinkBigCachePlus loadFromDisk SKIP partition [%s]: no FromByte registered", entry.Partition)
			continue
		}
		data, err := fromByte(entry.Value)
		if err != nil {
			// 单条解码失败跳过，不中断整次加载。
			log.Error("ThinkBigCachePlus loadFromDisk SKIP entry [%s/%s]: %s", entry.Partition, entry.Key, err.Error())
			continue
		}
		p.m_mapData.Store(entry.Key, data)
	}

	return nil
}

// SaveToDisk 把各分区数据逐条流式编码写盘，允许调用者手动落盘。
// 每个 partition 内是同一类数据，共用一个 ToByte；若分区未注册 ToByte 则跳过该分区。
func (this *ThinkBigCachePlus) SaveToDisk() error {
	this.m_lockFile.Lock()
	defer this.m_lockFile.Unlock()

	tmp := this.m_szFileName + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// 逐条流式编码写盘：任意时刻内存里只有一条记录 + 1MB 缓冲，避免整表拷进 buffer 造成的瞬时分配与 GC 尖峰。
	w := bufio.NewWriterSize(f, 1<<20)
	enc := gob.NewEncoder(w)

	var entry thinkBigCachePlusDiskEntry
	var encErr error
	nSaved := 0

	this.m_mapPartition.Range(func(k, v any) bool {
		szPartition := k.(string)
		p := v.(*ThinkBigCachePlusPartition)
		toByte := p.m_funcToByte
		if nil == toByte {
			// 分区未注册 ToByte，跳过整个分区。
			log.Error("ThinkBigCachePlus saveToDisk SKIP partition [%s]: no ToByte registered", szPartition)
			return true
		}

		p.m_mapData.Range(func(kk, vv any) bool {
			b, err := toByte(vv)
			if err != nil {
				// 单条编码失败跳过，不中断整次落盘。
				log.Error("ThinkBigCachePlus saveToDisk SKIP entry [%s/%s]: %s", szPartition, kk.(string), err.Error())
				return true
			}
			entry.Partition = szPartition
			entry.Key = kk.(string)
			entry.Value = b
			if err := enc.Encode(&entry); err != nil {
				log.Error(err.Error())
				encErr = err
				return false
			}

			nSaved++
			if 0 == nSaved%50 {
				time.Sleep(time.Millisecond)
			}
			return true
		})
		return nil == encErr
	})

	if encErr != nil {
		_ = f.Close()
		return encErr
	}

	if err := w.Flush(); err != nil {
		log.Error(err.Error())
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		log.Error(err.Error())
		return err
	}

	return os.Rename(tmp, this.m_szFileName)
}

func (this *ThinkBigCachePlus) emitUpdate(nType ThinkBigCachePlusUpdateType) {

	switch nType {
	case THINK_BIGCACHE_PLUS_UPDATE_1_MIN:
		for _, pFunc := range this.m_lst1MinListener {
			go pFunc()
		}
	case THINK_BIGCACHE_PLUS_UPDATE_5_MIN:
		for _, pFunc := range this.m_lst5MinListener {
			go pFunc()
		}
	case THINK_BIGCACHE_PLUS_UPDATE_10_MIN:
		for _, pFunc := range this.m_lst10MinListener {
			go pFunc()
		}
	case THINK_BIGCACHE_PLUS_UPDATE_30_MIN:
		for _, pFunc := range this.m_lst30MinListener {
			go pFunc()
		}
	case THINK_BIGCACHE_PLUS_UPDATE_1_HOUR:
		for _, pFunc := range this.m_lst1HourListener {
			go pFunc()
		}
	case THINK_BIGCACHE_PLUS_UPDATE_6_HOUR:
		for _, pFunc := range this.m_lst6HourListener {
			go pFunc()
		}
	case THINK_BIGCACHE_PLUS_UPDATE_12_HOUR:
		for _, pFunc := range this.m_lst12HourListener {
			go pFunc()
		}
	case THINK_BIGCACHE_PLUS_UPDATE_1_DAY:
		for _, pFunc := range this.m_lst1DayListener {
			go pFunc()
		}
	}

	if this.m_nSaveToDiskType == nType {
		go func() {
			time.Sleep(60 * time.Second)
			_ = this.SaveToDisk()
		}()
	}
}

func (this *ThinkBigCachePlus) initCron() error {
	this.m_pCronJobs = gocron.NewScheduler(time.Local)

	this.m_lst1MinListener = make([]OnThinkBigCachePlusUpdate, 0)
	this.m_lst5MinListener = make([]OnThinkBigCachePlusUpdate, 0)
	this.m_lst10MinListener = make([]OnThinkBigCachePlusUpdate, 0)
	this.m_lst30MinListener = make([]OnThinkBigCachePlusUpdate, 0)
	this.m_lst1HourListener = make([]OnThinkBigCachePlusUpdate, 0)
	this.m_lst6HourListener = make([]OnThinkBigCachePlusUpdate, 0)
	this.m_lst12HourListener = make([]OnThinkBigCachePlusUpdate, 0)
	this.m_lst1DayListener = make([]OnThinkBigCachePlusUpdate, 0)

	_, _ = this.m_pCronJobs.Cron("* * * * *").Do(func() {
		this.emitUpdate(THINK_BIGCACHE_PLUS_UPDATE_1_MIN)
	})

	_, _ = this.m_pCronJobs.Cron("*/5 * * * *").Do(func() {
		this.emitUpdate(THINK_BIGCACHE_PLUS_UPDATE_5_MIN)
		// 仅在进程启动后落盘一次，之后的周期落盘由整点任务负责。
		if this.m_bSavedToDisk.CompareAndSwap(false, true) {
			time.Sleep(60 * time.Second)
			_ = this.SaveToDisk()
		}
	})

	_, _ = this.m_pCronJobs.Cron("*/10 * * * *").Do(func() {
		this.emitUpdate(THINK_BIGCACHE_PLUS_UPDATE_10_MIN)
	})

	_, _ = this.m_pCronJobs.Cron("*/30 * * * *").Do(func() {
		this.emitUpdate(THINK_BIGCACHE_PLUS_UPDATE_30_MIN)
	})

	_, _ = this.m_pCronJobs.Cron("0 * * * *").Do(func() {
		this.emitUpdate(THINK_BIGCACHE_PLUS_UPDATE_1_HOUR)
	})

	_, _ = this.m_pCronJobs.Cron("0 0,6,12,18 * * *").Do(func() {
		this.emitUpdate(THINK_BIGCACHE_PLUS_UPDATE_6_HOUR)
	})

	_, _ = this.m_pCronJobs.Cron("0 0,12 * * *").Do(func() {
		this.emitUpdate(THINK_BIGCACHE_PLUS_UPDATE_12_HOUR)
	})

	_, _ = this.m_pCronJobs.Cron("0 0 * * *").Do(func() {
		this.emitUpdate(THINK_BIGCACHE_PLUS_UPDATE_1_DAY)
	})

	this.m_pCronJobs.StartAsync()

	return nil
}

func (this *ThinkBigCachePlus) Start() error {
	return this.StartEx(THINK_BIGCACHE_PLUS_UPDATE_1_DAY)
}

func (this *ThinkBigCachePlus) StartEx(nSaveToDiskType ThinkBigCachePlusUpdateType) error {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()

	if this.m_bStarted.Load() {
		return nil
	}

	this.m_nSaveToDiskType = nSaveToDiskType
	var err error = nil

	err = this.loadFromDisk()
	if err != nil {
		goto err_ret
	}

	err = this.initCron()
	if err != nil {
		goto err_ret
	}

	this.m_bStarted.Store(true)
	log.Info("ThinkBigCachePlus started successfully")

err_ret:
	return err
}

func (this *ThinkBigCachePlus) AddUpdateListener(nType ThinkBigCachePlusUpdateType, pFunc OnThinkBigCachePlusUpdate) {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()
	switch nType {
	case THINK_BIGCACHE_PLUS_UPDATE_1_MIN:
		this.m_lst1MinListener = append(this.m_lst1MinListener, pFunc)
	case THINK_BIGCACHE_PLUS_UPDATE_5_MIN:
		this.m_lst5MinListener = append(this.m_lst5MinListener, pFunc)
	case THINK_BIGCACHE_PLUS_UPDATE_10_MIN:
		this.m_lst10MinListener = append(this.m_lst10MinListener, pFunc)
	case THINK_BIGCACHE_PLUS_UPDATE_30_MIN:
		this.m_lst30MinListener = append(this.m_lst30MinListener, pFunc)
	case THINK_BIGCACHE_PLUS_UPDATE_1_HOUR:
		this.m_lst1HourListener = append(this.m_lst1HourListener, pFunc)
	case THINK_BIGCACHE_PLUS_UPDATE_6_HOUR:
		this.m_lst6HourListener = append(this.m_lst6HourListener, pFunc)
	case THINK_BIGCACHE_PLUS_UPDATE_12_HOUR:
		this.m_lst12HourListener = append(this.m_lst12HourListener, pFunc)
	case THINK_BIGCACHE_PLUS_UPDATE_1_DAY:
		this.m_lst1DayListener = append(this.m_lst1DayListener, pFunc)
	}
}

// Set 存入数据。内存中直接保存实际 struct(any)，读取时零解码、无额外 GC 压力。
func (this *ThinkBigCachePlus) Set(szPartition, szKey string, data any) error {
	p := this.getOrCreatePartition(szPartition)
	p.m_mapData.Store(szKey, data)
	return nil
}

// Get 获取数据，直接返回内存中的实际 struct(any)。
func (this *ThinkBigCachePlus) Get(szPartition, szKey string) (any, error) {
	v, ok := this.m_mapPartition.Load(szPartition)
	if !ok {
		return nil, errors.New("ThinkBigCachePlus partition not found: " + szPartition)
	}
	p := v.(*ThinkBigCachePlusPartition)
	if data, ok := p.m_mapData.Load(szKey); ok {
		return data, nil
	}
	return nil, errors.New("ThinkBigCachePlus key not found: " + szKey)
}

// RegToByteFunction 为指定分区一次性注册编解码对。
// 每个 partition 内是同一类数据，故编解码方法一致：ToByte 用于 saveToDisk，FromByte 用于 loadFromDisk。
// 应在 Start 之前调用，以便 loadFromDisk 能正确把字节还原成 struct。
func (this *ThinkBigCachePlus) RegByteFunction(szPartition string, pToByte ThinkBigCachePlusToByte, pFromByte ThinkBigCachePlusFromByte) {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()

	p := this.getOrCreatePartition(szPartition)
	p.m_funcToByte = pToByte
	p.m_funcFromByte = pFromByte
}
