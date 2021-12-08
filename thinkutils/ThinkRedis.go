package thinkutils

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
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
