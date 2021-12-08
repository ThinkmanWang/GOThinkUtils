package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"context"
	"sync"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func redisTest(wg *sync.WaitGroup) {
	defer wg.Done()

	rdb := thinkutils.ThinkRedis.QuickConn()
	rdb.Del(context.Background(), "FXXK")

	err := rdb.Set(context.Background(), "FXXK", thinkutils.DateTime.CurDatetime(), 0).Err()
	if err != nil {
		return
	}

	szVal, err := rdb.Get(context.Background(), "FXXK").Result()
	if err != nil {
		return
	}

	log.Info(szVal)
}

func lockTest(wg *sync.WaitGroup) {
	defer wg.Done()

	rdb := thinkutils.ThinkRedis.QuickConn()
	szVal := thinkutils.ThinkRedis.Lock(rdb, "mylock", 10, 60)
	log.Info(szVal)
	thinkutils.ThinkRedis.ReleaseLock(rdb, "mylock", szVal)
}

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 1024; i++ {
		wg.Add(1)
		go redisTest(&wg)
	}

	wg.Add(1)
	go lockTest(&wg)

	wg.Wait()
	//time.Sleep(10 * time.Second)
	//basicQueryJSON()
}
