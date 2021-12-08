package thinkutils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gopkg.in/ini.v1"
	"sync"
	"time"
)

type thinkredis struct {
}

var (
	g_lockRedis    sync.Mutex
	g_mapRedisConn map[string]*redis.Client
)

func (this thinkredis) makeConn(szHost string,
	nPort int,
	szPwd string,
	nDb int,
	nMaxConn int) *redis.Client {
	defer g_lockRedis.Unlock()
	g_lockRedis.Lock()

	szKey := fmt.Sprintf("%s@(%s:%d)/%d", szPwd, szHost, nPort, nDb)
	rdb := g_mapRedisConn[szKey]
	if nil == rdb {
		rdb = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", szHost, nPort),
			Password:     szPwd,
			DB:           nDb,
			MinIdleConns: 2,
			PoolSize:     nMaxConn,
		})

		g_mapRedisConn[szKey] = rdb
	}

	return rdb
}

func (this thinkredis) Conn(szHost string,
	nPort int,
	szPwd string,
	nDb int,
	nMaxConn int) *redis.Client {

	if nil == g_mapRedisConn {
		g_mapRedisConn = make(map[string]*redis.Client)
	}

	szKey := fmt.Sprintf("%s@(%s:%d)/%d", szPwd, szHost, nPort, nDb)
	rdb := g_mapRedisConn[szKey]
	if nil == rdb {
		rdb = this.makeConn(szHost, nPort, szPwd, nDb, nMaxConn)
	}

	//log.Info("%p %p", g_mapRedisConn, rdb)
	return rdb
}

func (this thinkredis) QuickConn() *redis.Client {
	cfg, err := ini.Load("app.ini")
	if err != nil {
		return this.Conn("127.0.0.1", 6379, "123456", 0, 32)
	}

	return this.Conn(cfg.Section("redis").Key("host").String(),
		cfg.Section("redis").Key("port").MustInt(),
		cfg.Section("redis").Key("password").String(),
		cfg.Section("redis").Key("db").MustInt(),
		cfg.Section("redis").Key("max_conn").MustInt())
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
