package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"fmt"
	"sync"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	nIndex := 0
	for {
		nIndex++
		szMsg := fmt.Sprintf("[%d] %s", nIndex, thinkutils.DateTime.CurDatetime())
		thinkutils.KafkaUtils.SendMsg("172.16.0.2:9092",
			"think-topic",
			[]byte(szMsg))

		log.Info("Send %s", szMsg)
		time.Sleep(time.Duration(10) * time.Microsecond)
	}

	wg.Wait()
}
