package thinkutils

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"reflect"
	gort "runtime"
	"sync"
	"time"
)

type FuncRefreshCache func() error

type CacheNode struct {
	m_nExpireAt int64
	m_pFunc     FuncRefreshCache
	m_pData     any
}

type ThinkMemCache struct {
	m_dictCache map[string]*CacheNode
	m_lock      sync.RWMutex
	m_bStarted  bool
}

var (
	g_pInstance *ThinkMemCache = &ThinkMemCache{
		m_dictCache: make(map[string]*CacheNode),
		m_bStarted:  false,
	}
)

func GetMemCacheInstance() *ThinkMemCache {
	return g_pInstance
}

func (this *ThinkMemCache) Start() {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()

	if this.m_bStarted {
		return
	}

	go func() {
		for {
			time.Sleep(60 * time.Second)

			GetMemCacheInstance().Refresh()
		}
	}()
}

func (this *ThinkMemCache) Set(szKey string, nExpire int64, pData any, pFunc FuncRefreshCache) {
	nExpireAt := int64(0)
	if nExpire > 0 {
		nExpireAt = DateTime.Timestamp() + nExpire
	}

	pNode := &CacheNode{
		m_nExpireAt: nExpireAt,
		m_pData:     pData,
		m_pFunc:     pFunc,
	}

	this.m_lock.Lock()
	this.m_dictCache[szKey] = pNode
	this.m_lock.Unlock()
}

func (this *ThinkMemCache) get(szKey string) *CacheNode {
	this.m_lock.RLock()
	defer this.m_lock.RUnlock()

	pNode, bExists := this.m_dictCache[szKey]
	if false == bExists {
		return nil
	}

	if nil == pNode {
		return nil
	}

	return pNode
}

func (this *ThinkMemCache) Get(szKey string) any {
	pNode := this.get(szKey)
	if nil == pNode {
		return nil
	}

	return pNode.m_pData
}

func (this *ThinkMemCache) keys() []string {
	this.m_lock.RLock()
	defer this.m_lock.RUnlock()

	lstRet := make([]string, 0)
	for k, _ := range this.m_dictCache {
		lstRet = append(lstRet, k)
	}

	return lstRet
}

func (this *ThinkMemCache) Refresh() {
	logger.Info(">>>>REFRESH ALL CACHE<<<<")
	defer logger.Info(">>>>FINISH REFRESH ALL CACHE<<<<")

	//logger.Info("Cache size: %d", len(this.m_dictCache))
	lstKeys := this.keys()
	if nil == lstKeys || len(lstKeys) <= 0 {
		logger.Info("NO CACHE")
		return
	}

	wg := sync.WaitGroup{}
	for _, szKey := range lstKeys {
		v := this.get(szKey)
		if nil == v {
			continue
		}

		if 0 == v.m_nExpireAt {
			logger.Info("[%s] => %s never expire", szKey, gort.FuncForPC(reflect.ValueOf(v.m_pFunc).Pointer()).Name())
		} else {
			logger.Info("[%s] => %s expire in %d s", szKey, gort.FuncForPC(reflect.ValueOf(v.m_pFunc).Pointer()).Name(), v.m_nExpireAt-DateTime.Timestamp())
		}

		if v.m_nExpireAt != 0 && DateTime.Timestamp() >= v.m_nExpireAt && v.m_pFunc != nil {
			wg.Add(1)
			go func(pFunc FuncRefreshCache) {
				defer wg.Done()

				nStart := DateTime.TimestampMs()
				if err := pFunc(); err != nil {
					logger.Warn(err.Error())
				}
				logger.Info(">>>>Call %s for cache, cost %d<<<<", gort.FuncForPC(reflect.ValueOf(pFunc).Pointer()).Name(), DateTime.TimestampMs()-nStart)
			}(v.m_pFunc)

		}
	}
	wg.Wait()
}
