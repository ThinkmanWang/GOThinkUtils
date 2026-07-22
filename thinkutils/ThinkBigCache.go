package thinkutils

import (
	"bytes"
	"context"
	"encoding/gob"
	"os"
	"sync"
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
	m_lock      sync.Mutex
	m_bStarted  bool
	m_pCronJobs *gocron.Scheduler
	m_pBigCache *bigcache.BigCache

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
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(c)
}

func (this *ThinkBigCache) saveToDisk(c *bigcache.BigCache, path string) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(c); err != nil {
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
		this.saveToDisk(this.m_pBigCache, "ThinkBigCache.data")
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

func (this *ThinkBigCache) Start() error {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()

	if this.m_bStarted {
		return nil
	}

	var err error = nil

	cfg := bigcache.Config{
		Shards:           1024,
		LifeWindow:       100 * 365 * 24 * time.Hour, // 逻辑永久有效
		CleanWindow:      0,                          // 关闭后台清理，库绝不主动删除数据
		MaxEntrySize:     4 * 1024 * 1024,
		HardMaxCacheSize: 0, // 无内存上限
		Verbose:          false,
	}

	this.m_pBigCache, _ = bigcache.New(context.Background(), cfg)
	err = this.loadFromDisk(this.m_pBigCache, "ThinkBigCache.data")
	if err != nil {
		goto err_ret
	}

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
	return this.m_pBigCache.Set(szKey, data)
}

func (this *ThinkBigCache) Get(szKey string) ([]byte, error) {
	if data, err := this.m_pBigCache.Get(szKey); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}
