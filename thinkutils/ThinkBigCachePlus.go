package thinkutils

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"os"
	"reflect"
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
	// m_reflectType 记录该分区注册时的元素类型，用于类型化句柄的一致性校验(一个分区一个类型)。
	m_reflectType reflect.Type
	// m_bSaveToDisk 标记该分区是否需要落盘，默认 true；为 false 时 SaveToDisk 跳过该分区。
	m_bSaveToDisk atomic.Bool
}

type ThinkBigCachePlus struct {
	m_bStarted atomic.Bool
	//m_bSavedToDisk atomic.Bool
	m_lock       sync.RWMutex
	m_lockFile   sync.Mutex
	m_szFileName string

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

		m_lst1MinListener:   make([]OnThinkBigCachePlusUpdate, 0),
		m_lst5MinListener:   make([]OnThinkBigCachePlusUpdate, 0),
		m_lst10MinListener:  make([]OnThinkBigCachePlusUpdate, 0),
		m_lst30MinListener:  make([]OnThinkBigCachePlusUpdate, 0),
		m_lst1HourListener:  make([]OnThinkBigCachePlusUpdate, 0),
		m_lst6HourListener:  make([]OnThinkBigCachePlusUpdate, 0),
		m_lst12HourListener: make([]OnThinkBigCachePlusUpdate, 0),
		m_lst1DayListener:   make([]OnThinkBigCachePlusUpdate, 0),
	}
)

func ThinkBigCachePlusInstance() *ThinkBigCachePlus {
	return g_pThinkBigCachePlusInstance
}

// getOrCreatePartition 返回指定分区(不存在则原子创建)。
// bSaveToDisk 仅在首次创建时用于初始化该分区是否需要落盘的标记(默认调用方传 true)。
// 存指针而非值：ThinkBigCachePlusPartition 内含 sync.Map，绝不能被拷贝。
func (this *ThinkBigCachePlus) getOrCreatePartition(szPartition string, bSaveToDisk bool) *ThinkBigCachePlusPartition {
	if v, ok := this.m_mapPartition.Load(szPartition); ok {
		return v.(*ThinkBigCachePlusPartition)
	}
	p := &ThinkBigCachePlusPartition{}
	p.m_bSaveToDisk.Store(bSaveToDisk)
	actual, _ := this.m_mapPartition.LoadOrStore(szPartition, p)
	return actual.(*ThinkBigCachePlusPartition)
}

// NeedSaveToDisk 控制指定分区是否需要落盘。默认所有分区都会落盘；
// 设为 false 后，SaveToDisk 会跳过该分区(内存数据保留，只是不写盘)。
func (this *ThinkBigCachePlus) NeedSaveToDisk(szPartition string, bNeed bool) {
	p := this.getOrCreatePartition(szPartition, bNeed)
	p.m_bSaveToDisk.Store(bNeed)
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

		p := this.getOrCreatePartition(entry.Partition, true)
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
		if !p.m_bSaveToDisk.Load() {
			// 该分区被标记为不落盘(NeedSaveToDisk=false)，整个分区跳过。
			return true
		}
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

	_, _ = this.m_pCronJobs.Cron("* * * * *").Do(func() {
		this.emitUpdate(THINK_BIGCACHE_PLUS_UPDATE_1_MIN)
	})

	_, _ = this.m_pCronJobs.Cron("*/5 * * * *").Do(func() {
		this.emitUpdate(THINK_BIGCACHE_PLUS_UPDATE_5_MIN)
		// 仅在进程启动后落盘一次，之后的周期落盘由整点任务负责。
		//if this.m_bSavedToDisk.CompareAndSwap(false, true) {
		//	time.Sleep(60 * time.Second)
		//	_ = this.SaveToDisk()
		//}
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

// ThinkBigCachePlusPartitionT 是某分区的类型化句柄，提供编译期类型安全的 Set/Get。
// 内存中直接保存实际 struct(T)，Get 只做一次类型断言、零解码、无额外 GC 压力。
type ThinkBigCachePlusPartitionT[T any] struct {
	m_pPartition *ThinkBigCachePlusPartition
	m_szName     string
}

// Set 存入数据。data 的类型在编译期锁定为 T，无法传错类型。
func (this *ThinkBigCachePlusPartitionT[T]) Set(szKey string, data T) error {
	this.m_pPartition.m_mapData.Store(szKey, data)
	return nil
}

// Get 直接返回具体类型 T，调用端无需 .(T) 断言。
func (this *ThinkBigCachePlusPartitionT[T]) Get(szKey string) (T, error) {
	var zero T
	data, ok := this.m_pPartition.m_mapData.Load(szKey)
	if !ok {
		return zero, errors.New("ThinkBigCachePlus key not found: " + szKey)
	}
	if v, ok := data.(T); ok {
		return v, nil
	}
	return zero, errors.New("ThinkBigCachePlus type mismatch in partition: " + this.m_szName)
}

// Has 判断指定 key 是否存在(不做类型断言，仅判断键存在与否)。
func (this *ThinkBigCachePlusPartitionT[T]) Has(szKey string) bool {
	_, ok := this.m_pPartition.m_mapData.Load(szKey)
	return ok
}

// Delete 删除指定 key。用于需要"全量替换/解封"语义的分区。
func (this *ThinkBigCachePlusPartitionT[T]) Delete(szKey string) {
	this.m_pPartition.m_mapData.Delete(szKey)
}

// Range 遍历分区内所有键值，回调返回 false 时停止遍历。
// 与 sync.Map.Range 语义一致：遍历期间的并发写不保证被看到，且不应假设快照一致性。
// 类型断言失败的条目会被跳过(理论上不会发生，一个分区一个类型)。
func (this *ThinkBigCachePlusPartitionT[T]) Range(f func(szKey string, data T) bool) {
	this.m_pPartition.m_mapData.Range(func(k, v any) bool {
		val, ok := v.(T)
		if !ok {
			return true
		}
		return f(k.(string), val)
	})
}

// ThinkBigCachePlusRegType 为分区注册"一个类型"，内部用 gob 自动生成编解码，返回类型化句柄。
// 一个分区一个类型；应在 Start 之前调用(通常放在各模块 init)。
// 注意：gob 只编码导出字段；每条记录独立编码，会各自携带一份类型描述(体积/CPU 有额外开销，大数据量时留意)。
func ThinkBigCachePlusRegType[T any](this *ThinkBigCachePlus, szPartition string) *ThinkBigCachePlusPartitionT[T] {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()

	p := this.getOrCreatePartition(szPartition, true)
	p.m_reflectType = reflect.TypeOf((*T)(nil)).Elem()
	p.m_funcToByte = func(v any) ([]byte, error) {
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(v); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	p.m_funcFromByte = func(b []byte) (any, error) {
		var out T
		if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&out); err != nil {
			return nil, err
		}
		return out, nil
	}
	return &ThinkBigCachePlusPartitionT[T]{m_pPartition: p, m_szName: szPartition}
}

// ThinkBigCachePlusPartitionOf 返回已通过 RegType 注册分区的类型化句柄，不改动已注册的编解码。
// 适用于在一处 RegType 注册、在别处按分区名再取类型化句柄使用的场景。
// 若该分区已记录过类型且与 T 不一致，则记录告警(一个分区一个类型)。
func ThinkBigCachePlusPartitionOf[T any](this *ThinkBigCachePlus, szPartition string) *ThinkBigCachePlusPartitionT[T] {
	p := this.getOrCreatePartition(szPartition, true)
	if t := reflect.TypeOf((*T)(nil)).Elem(); p.m_reflectType != nil && p.m_reflectType != t {
		log.Error("ThinkBigCachePlus partition [%s] type mismatch: registered %s, requested %s", szPartition, p.m_reflectType.String(), t.String())
	}
	return &ThinkBigCachePlusPartitionT[T]{m_pPartition: p, m_szName: szPartition}
}
