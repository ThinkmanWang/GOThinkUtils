package main

import (
	"fmt"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"sync"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
	//bus EventBus.Bus        = EventBus.New()
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
	thinkutils.ThinkEventBus.Publish("main:message", "FXXK")

	pData := &MyStruct{Name: "123", Desc: "456"}
	log.Info("%p", pData)
	thinkutils.ThinkEventBus.Publish("main:otherMessage", pData)
	thinkutils.ThinkEventBus.Publish("main:multiParams", "123", "456")
	thinkutils.ThinkEventBus.Publish("main:123456", "12356")
}

func main() {
	thinkutils.ThinkEventBus.SubscribeAsync("main:message", onMsg, false)
	thinkutils.ThinkEventBus.SubscribeAsync("main:otherMessage", onStructMsg, false)
	thinkutils.ThinkEventBus.SubscribeAsync("main:multiParams", onMultiParams, false)

	go func() {
		publish()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	//bus.Unsubscribe("main:message", onMsg)
}
