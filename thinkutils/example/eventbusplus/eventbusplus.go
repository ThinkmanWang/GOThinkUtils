package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"fmt"
	"github.com/asaskevich/EventBus"
	"runtime"
	"sync"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
	bus EventBus.Bus        = EventBus.New()
)

type MyListener struct {
	Id uint32
}

func (this *MyListener) Init() {
	szMsg := fmt.Sprintf("sendTo:%d", this.Id)
	bus.SubscribeAsync(szMsg, this.sendTo, false)
}

func (this *MyListener) sendTo(msg string) {
	log.Info("[%p] %d -> %s", this, this.Id, msg)
}

func publish() {
	bus.Publish("sendTo:1", thinkutils.DateTime.CurDatetime())

	time.Sleep(2 * time.Second)
	bus.Publish("sendTo:2", thinkutils.DateTime.CurDatetime())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	listener1 := MyListener{Id: 1}
	log.Info("%p", &listener1)
	listener1.Init()

	listener2 := MyListener{Id: 2}
	log.Info("%p", &listener2)
	listener2.Init()

	go publish()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
