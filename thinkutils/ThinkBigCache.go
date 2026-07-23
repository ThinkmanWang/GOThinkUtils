package thinkutils

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/go-co-op/gocron"
)

type OnThinkBigCacheUpdate func()

type ThinkBigCacheUpdateType int

const (
	UPDATE_1_MIN ThinkBigCacheUpdateType = iota
	UPDATE_5_MIN
	UPDATE_10_MIN
	UPDATE_30_MIN
	UPDATE_1_HOUR
	UPDATE_6_HOUR
	UPDATE_12_HOUR
	UPDATE_1_DAY
)

type ThinkBigCache struct {
	m_lock     sync.RWMutex
	m_lockFile sync.Mutex
	// m_swapMu 串行化写端(Set)与 refreshMemory 重建：
	// 重建期间禁止并发 Set，避免写入落到即将被丢弃的旧 cache 上导致丢数据。
	// 读端(Get)不参与此锁，通过 m_pBigCache 原子读取当前实例，全程无锁。
	m_swapMu sync.Mutex

	m_bStarted bool
	m_config   bigcache.Config

	m_pCronJobs *gocron.Scheduler
	// m_pBigCache 采用原子指针：读端一次原子读拿到当前实例，写端(refreshMemory)原子替换整个实例。
	m_pBigCache atomic.Pointer[bigcache.BigCache]

	m_lst1MinListener   []OnThinkBigCacheUpdate
	m_lst5MinListener   []OnThinkBigCacheUpdate
	m_lst10MinListener  []OnThinkBigCacheUpdate
	m_lst30MinListener  []OnThinkBigCacheUpdate
	m_lst1HourListener  []OnThinkBigCacheUpdate
	m_lst6HourListener  []OnThinkBigCacheUpdate
	m_lst12HourListener []OnThinkBigCacheUpdate
	m_lst1DayListener   []OnThinkBigCacheUpdate
}

var (
	g_pThinkBigCacheInstance *ThinkBigCache = &ThinkBigCache{
		m_bStarted: false,
	}
)

func ThinkBigCacheInstance() *ThinkBigCache {
	return g_pThinkBigCacheInstance
}

func (this *ThinkBigCache) loadFromDisk(c *bigcache.BigCache, path string) error {
	this.m_lockFile.Lock()
	defer this.m_lockFile.Unlock()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	entries := make(map[string][]byte)
	if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&entries); err != nil {
		log.Error(err.Error())
		return err
	}

	for key, value := range entries {
		if err := c.Set(key, value); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}

func (this *ThinkBigCache) saveToDisk(c *bigcache.BigCache, path string) error {
	this.m_lockFile.Lock()
	defer this.m_lockFile.Unlock()

	entries := make(map[string][]byte)

	iter := c.Iterator()
	for iter.SetNext() {
		info, err := iter.Value()
		if err != nil {
			log.Error(err.Error())
			return err
		}
		entries[info.Key()] = info.Value()
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(entries); err != nil {
		log.Error(err.Error())
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0644); err != nil {
		log.Error(err.Error())
		return err
	}
	return os.Rename(tmp, path)
}

func (this *ThinkBigCache) emitUpdate(nType ThinkBigCacheUpdateType) {
	this.m_lock.RLock()
	defer this.m_lock.RUnlock()

	switch nType {
	case UPDATE_1_MIN:
		for _, pFunc := range this.m_lst1MinListener {
			go pFunc()
		}
	case UPDATE_5_MIN:
		for _, pFunc := range this.m_lst5MinListener {
			go pFunc()
		}
	case UPDATE_10_MIN:
		for _, pFunc := range this.m_lst10MinListener {
			go pFunc()
		}
	case UPDATE_30_MIN:
		for _, pFunc := range this.m_lst30MinListener {
			go pFunc()
		}
	case UPDATE_1_HOUR:
		for _, pFunc := range this.m_lst1HourListener {
			go pFunc()
		}
	case UPDATE_6_HOUR:
		for _, pFunc := range this.m_lst6HourListener {
			go pFunc()
		}
	case UPDATE_12_HOUR:
		for _, pFunc := range this.m_lst12HourListener {
			go pFunc()
		}
	case UPDATE_1_DAY:
		for _, pFunc := range this.m_lst1DayListener {
			go pFunc()
		}
	}
}

func (this *ThinkBigCache) initCron() error {
	this.m_pCronJobs = gocron.NewScheduler(time.Local)

	this.m_lst1MinListener = make([]OnThinkBigCacheUpdate, 0)
	this.m_lst5MinListener = make([]OnThinkBigCacheUpdate, 0)
	this.m_lst10MinListener = make([]OnThinkBigCacheUpdate, 0)
	this.m_lst30MinListener = make([]OnThinkBigCacheUpdate, 0)
	this.m_lst1HourListener = make([]OnThinkBigCacheUpdate, 0)
	this.m_lst6HourListener = make([]OnThinkBigCacheUpdate, 0)
	this.m_lst12HourListener = make([]OnThinkBigCacheUpdate, 0)
	this.m_lst1DayListener = make([]OnThinkBigCacheUpdate, 0)

	_, _ = this.m_pCronJobs.Cron("* * * * *").Do(func() {
		this.emitUpdate(UPDATE_1_MIN)
	})

	_, _ = this.m_pCronJobs.Cron("*/5 * * * *").Do(func() {
		this.saveToDisk(this.m_pBigCache.Load(), "ThinkBigCache.data")
		this.emitUpdate(UPDATE_5_MIN)
	})

	_, _ = this.m_pCronJobs.Cron("*/10 * * * *").Do(func() {
		this.emitUpdate(UPDATE_10_MIN)
	})

	_, _ = this.m_pCronJobs.Cron("*/30 * * * *").Do(func() {
		this.emitUpdate(UPDATE_30_MIN)
	})

	_, _ = this.m_pCronJobs.Cron("0 * * * *").Do(func() {
		this.emitUpdate(UPDATE_1_HOUR)
		_ = this.refreshMemory()
		runtime.GC()
	})

	_, _ = this.m_pCronJobs.Cron("0 0,6,12,18 * * *").Do(func() {
		this.emitUpdate(UPDATE_6_HOUR)
	})

	_, _ = this.m_pCronJobs.Cron("0 0,12 * * *").Do(func() {
		this.emitUpdate(UPDATE_12_HOUR)
	})

	_, _ = this.m_pCronJobs.Cron("0 0 * * *").Do(func() {
		this.emitUpdate(UPDATE_1_DAY)
	})

	this.m_pCronJobs.StartAsync()

	return nil
}

// refreshMemory 重建整个 BigCache 以回收覆盖写产生的死空间(dead space)。
// BigCache 对已存在 key 的 Set 不会原地更新，而是追加新数据并将旧条置无效；
// 在 CleanWindow=0 + HardMaxCacheSize=0 下旧条字节永不回收，内存会随刷新次数单调上涨。
// 此处新建一个 cache，将旧 cache 中的活数据全量复制过去，再原子替换实例，从而丢弃死空间。
func (this *ThinkBigCache) refreshMemory() error {
	// 与 Set 互斥：重建期间不允许并发写入，保证旧 cache 在复制过程中不再变化，避免丢写。
	this.m_swapMu.Lock()
	defer this.m_swapMu.Unlock()

	old := this.m_pBigCache.Load()
	if nil == old {
		return nil
	}

	nStart := DateTime.TimestampMs()
	log.Info("ThinkBigCache refreshMemory START")

	fresh, err := bigcache.New(context.Background(), this.m_config)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	nCopied := 0
	iter := old.Iterator()
	for iter.SetNext() {
		info, err := iter.Value()
		if err != nil {
			// 单条取值失败跳过，不中断整次重建。
			continue
		}
		if err := fresh.Set(info.Key(), info.Value()); err != nil {
			log.Error(err.Error())
			continue
		}
		nCopied++
	}

	// 原子发布新实例；旧实例丢引用后由 GC 回收。
	// CleanWindow=0 无后台 goroutine，无需显式 Close，且不 Close 可避免影响刚 Load 到旧实例、随后才 Get 的请求。
	this.m_pBigCache.Store(fresh)
	go func(o *bigcache.BigCache) {
		time.Sleep(30 * time.Second) // 给在途请求留出用完旧指针的时间
		_ = o.Close()
	}(old)
	log.Info("ThinkBigCache refreshMemory FINISH copied %d entries, cost %d ms", nCopied, DateTime.TimestampMs()-nStart)

	return nil
}

func (this *ThinkBigCache) Start() error {
	cfg := bigcache.Config{
		Shards:           1024,
		LifeWindow:       100 * 365 * 24 * time.Hour, // 逻辑永久有效
		CleanWindow:      0,                          // 关闭后台清理，库绝不主动删除数据
		MaxEntrySize:     64 * 1024,
		HardMaxCacheSize: 0, // 无内存上限
		Verbose:          false,
	}

	return this.StartEx(cfg)
}

func (this *ThinkBigCache) StartEx(cfg bigcache.Config) error {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()

	if this.m_bStarted {
		return nil
	}

	var err error = nil
	var pBigCache *bigcache.BigCache = nil

	this.m_config = cfg

	pBigCache, err = bigcache.New(context.Background(), this.m_config)
	if err != nil {
		goto err_ret
	}
	this.m_pBigCache.Store(pBigCache)

	err = this.loadFromDisk(pBigCache, "ThinkBigCache.data")
	if err != nil {
		goto err_ret
	}
	_ = this.Set("Hello", []byte("Hello World"))
	_ = this.saveToDisk(pBigCache, "ThinkBigCache.data")

	err = this.initCron()
	if err != nil {
		goto err_ret
	}

	this.m_bStarted = true
	log.Info("ThinkBigCache started successfully")

err_ret:
	return err
}

func (this *ThinkBigCache) AddUpdateListener(nType ThinkBigCacheUpdateType, pFunc OnThinkBigCacheUpdate) {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()
	switch nType {
	case UPDATE_1_MIN:
		this.m_lst1MinListener = append(this.m_lst1MinListener, pFunc)
	case UPDATE_5_MIN:
		this.m_lst5MinListener = append(this.m_lst5MinListener, pFunc)
	case UPDATE_10_MIN:
		this.m_lst10MinListener = append(this.m_lst10MinListener, pFunc)
	case UPDATE_30_MIN:
		this.m_lst30MinListener = append(this.m_lst30MinListener, pFunc)
	case UPDATE_1_HOUR:
		this.m_lst1HourListener = append(this.m_lst1HourListener, pFunc)
	case UPDATE_6_HOUR:
		this.m_lst6HourListener = append(this.m_lst6HourListener, pFunc)
	case UPDATE_12_HOUR:
		this.m_lst12HourListener = append(this.m_lst12HourListener, pFunc)
	case UPDATE_1_DAY:
		this.m_lst1DayListener = append(this.m_lst1DayListener, pFunc)
	}
}

func (this *ThinkBigCache) Set(szKey string, data []byte) error {
	// 与 refreshMemory 互斥：避免写入落到重建期间即将被丢弃的旧 cache 上导致丢数据。
	// Set 仅在配置刷新时调用，不在请求热路径，加锁开销可忽略。
	this.m_swapMu.Lock()
	defer this.m_swapMu.Unlock()

	pBigCache := this.m_pBigCache.Load()
	if nil == pBigCache {
		return errors.New("ThinkBigCache is nil")
	}

	return pBigCache.Set(szKey, data)
}

func (this *ThinkBigCache) Get(szKey string) ([]byte, error) {
	// 读端热路径：一次原子读拿到当前实例，全程无锁。
	pBigCache := this.m_pBigCache.Load()
	if nil == pBigCache {
		return nil, errors.New("ThinkBigCache is nil")
	}

	if data, err := pBigCache.Get(szKey); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}
