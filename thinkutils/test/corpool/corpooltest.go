package main

import (
	"GOThinkUtils/thinkutils/logger"
	"github.com/panjf2000/ants"
	"sync"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	wg := sync.WaitGroup{}
	pool, _ := ants.NewPool(10)

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			pool.Submit(func() {
				log.Info("FXXXXXXK")
				time.Sleep(time.Second)
				wg.Done()
			})
		}()
	}
	defer pool.Release()

	pool1, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		log.Info("FXXK %d", i)
		time.Sleep(time.Second)
		wg.Done()
	})
	defer pool1.Release()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		pool1.Invoke(i)
	}

	log.Info("Wait for pool finish")
	wg.Wait()
}
