package main

import (
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"fmt"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {

	nIndex := 0
	for i := 0; i < 10; i++ {
		nIndex++
		szMsg := fmt.Sprintf("[%d] %s", nIndex, thinkutils.DateTime.CurDatetime())
		thinkutils.KafkaUtils.SendMsg("172.16.0.2:9092",
			"think-topic",
			[]byte(szMsg))

		log.Info("Send %s", szMsg)
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

	time.Sleep(10 * time.Second)
}
