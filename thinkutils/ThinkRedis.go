package thinkutils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type thinkredis struct {
	m_lock    sync.Mutex
	m_mapConn map[string]*redis.Client
}

func (this thinkredis) Conn(szHost string,
	nPort int,
	szPwd string,
	nDb int,
	nMaxConn int) *redis.Client {
	defer this.m_lock.Unlock()
	this.m_lock.Lock()

	if nil == this.m_mapConn {
		this.m_mapConn = make(map[string]*redis.Client)
	}

	szKey := fmt.Sprintf("%s@(%s:%d)/%d", szPwd, szHost, nPort, nDb)
	rdb := this.m_mapConn[szKey]
	if rdb != nil {
		return rdb
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", szHost, nPort),
		Password:     szPwd,
		DB:           nDb,
		MinIdleConns: 2,
		PoolSize:     nMaxConn,
	})

	this.m_mapConn[szKey] = rdb
	return rdb
}

func (this thinkredis) QuickConn() *redis.Client {
	return this.Conn("172.16.0.2", 6379, "Ab123145", 0, 32)
}

func (this thinkredis) Lock(rDB *redis.Client, szName string, nAcquireTimeout int32, nLockTimeout int32) string {
	if nil == rDB {
		return ""
	}

	if StringUtils.IsEmpty(szName) {
		return ""
	}

	szLockName := fmt.Sprintf("lock:%s", szName)
	szVal := UUIDUtils.New()
	nEndTime := DateTime.Timestamp() + int64(nAcquireTimeout)

	for true {
		err := rDB.SetNX(context.Background(), szLockName, szVal, time.Duration(nLockTimeout)*time.Second).Err()
		if nil == err {
			return szVal
		}

		if DateTime.Timestamp() >= nEndTime {
			break
		}
	}

	return ""
}

func (this thinkredis) ReleaseLock(rDB *redis.Client, szName string, szVal string) {
	if nil == rDB {
		return
	}

	if StringUtils.IsEmpty(szName) {
		return
	}

	szLockName := fmt.Sprintf("lock:%s", szName)
	val, err := rDB.Get(context.Background(), szLockName).Result()
	if err != nil {
		return
	}

	if szVal == val {
		rDB.Del(context.Background(), szLockName)
	}
}
