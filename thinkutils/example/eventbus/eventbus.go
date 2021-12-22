package main

import (
	"GOThinkUtils/thinkutils/logger"
	"fmt"
	"github.com/asaskevich/EventBus"
	"sync"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
	bus EventBus.Bus        = EventBus.New()
)

func onMsg(szMsg string) {
	time.Sleep(5 * time.Second)
	log.Info(szMsg)
}

type MyStruct struct {
	Name string
	Desc string
}

func onStructMsg(data *MyStruct) {
	log.Info("%p", data)

	time.Sleep(5 * time.Second)
	fmt.Println(*data)
}

func onMultiParams(szName string, szDesc string) {
	time.Sleep(5 * time.Second)
	log.Info("%s %s", szName, szDesc)
}

func publish() {
	bus.Publish("main:message", "FXXK")

	pData := &MyStruct{Name: "123", Desc: "456"}
	log.Info("%p", pData)
	bus.Publish("main:otherMessage", pData)
	bus.Publish("main:multiParams", "123", "456")
	bus.Publish("main:123456", "12356")
}

func main() {
	bus.SubscribeAsync("main:message", onMsg, false)
	bus.SubscribeAsync("main:otherMessage", onStructMsg, false)
	bus.SubscribeAsync("main:multiParams", onMultiParams, false)

	go func() {
		publish()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	//bus.Unsubscribe("main:message", onMsg)
}
