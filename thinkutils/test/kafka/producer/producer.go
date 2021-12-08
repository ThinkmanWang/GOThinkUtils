package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"sync"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		thinkutils.KafkaUtils.SendMsg("172.16.0.2:9092",
			"think-topic",
			[]byte(thinkutils.DateTime.CurDatetime()))

		time.Sleep(time.Duration(1) * time.Second)
	}

	wg.Wait()
}
