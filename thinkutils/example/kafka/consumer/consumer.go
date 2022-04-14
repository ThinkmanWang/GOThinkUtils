package main

import (
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"github.com/segmentio/kafka-go"
	"sync"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	thinkutils.KafkaUtils.StartConsumer("172.16.0.2:9092,172.16.0.2:9093,172.16.0.2:9094",
		"think-topic",
		"thinkgo",
		func(message kafka.Message) {
			log.Info(thinkutils.StringUtils.BytesToString(message.Value))
		})

	wg.Wait()
}
