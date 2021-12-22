package main

import (
	"GOThinkUtils/thinkutils/logger"
	"runtime"
	"sync"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func cortest(nVal int) {
	//log.Info("%d", nVal)
	go func() {
		log.Info("%d", nVal)
	}()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{}
	wg.Add(1)

	for i := 0; i < 10; i++ {
		cortest(i)
	}

	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		go func() {
			log.Info("%d", i)
		}()
	}

	wg.Wait()
}
