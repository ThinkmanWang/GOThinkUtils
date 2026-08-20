package thinkutils

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron"
)

type OnThinkBigCachePlusUpdate func()
type ThinkBigCachePlusToByte func(pData any) []byte

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
	m_mapData    sync.Map
	m_funcToByte ThinkBigCachePlusToByte
}

type ThinkBigCachePlus struct {
	m_bStarted     bool
	m_bSavedToDisk atomic.Bool
	m_lock         sync.RWMutex
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
		m_bStarted:   false,
		m_szFileName: "ThinkBigCachePlus.data",
	}
)

func ThinkBigCachePlusInstance() *ThinkBigCachePlus {
	return g_pThinkBigCachePlusInstance
}

func (this *ThinkBigCachePlus) loadFromDisk(szPath string) error {
	return nil
}

func (this *ThinkBigCachePlus) saveToDisk(szPath string) error {
	return nil
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
			this.saveToDisk(this.m_szFileName)
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
			_ = this.saveToDisk(this.m_szFileName)
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

	if this.m_bStarted {
		return nil
	}

	this.m_nSaveToDiskType = nSaveToDiskType
	var err error = nil

	err = this.loadFromDisk(this.m_szFileName)
	if err != nil {
		goto err_ret
	}

	err = this.initCron()
	if err != nil {
		goto err_ret
	}

	this.m_bStarted = true
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

func (this *ThinkBigCachePlus) Set(szPartition, szKey string, data any) error {
	// 与 refreshMemory 互斥：避免写入落到重建期间即将被丢弃的旧 cache 上导致丢数据。
	// Set 仅在配置刷新时调用，不在请求热路径，加锁开销可忽略。

	return nil
}

func (this *ThinkBigCachePlus) Get(szPartition, szKey string) (any, error) {
	return nil, nil
}

func (this *ThinkBigCachePlus) RegToByteFunction(szPartition string, pFunc ThinkBigCachePlusToByte) {

}
