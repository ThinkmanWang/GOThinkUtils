package main

import (
	"GOThinkUtils/thinkutils/logger"
	"github.com/asaskevich/EventBus"
	"sync"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
	bus EventBus.Bus        = EventBus.New()
)

func onMsg(szMsg string) {
	log.Info(szMsg)
}

func publish() {
	bus.Publish("main:message", "FXXK")
}

func main() {
	bus.Subscribe("main:message", onMsg)

	go func() {
		publish()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	//bus.Unsubscribe("main:message", onMsg)
}
