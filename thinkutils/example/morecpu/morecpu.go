package main

import (
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"fmt"
	"runtime"
	"sync"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func task(wg *sync.WaitGroup) {
	defer wg.Done()

	var nSum uint64 = 0
	for i := 1; i < 9999999999; i++ {
		nSum += uint64(i)
	}

	fmt.Println(nSum)
}

func main() {

	//runtime.GOMAXPROCS(1)
	//wg := sync.WaitGroup{}
	//wg.Add(20)
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		fmt.Println("A: ", i)
	//		wg.Done()
	//	}()
	//}
	//for i := 0; i < 10; i++ {
	//	go func(i int) {
	//		fmt.Println("B: ", i)
	//		wg.Done()
	//	}(i)
	//}
	//wg.Wait()

	runtime.GOMAXPROCS(runtime.NumCPU())
	//runtime.GOMAXPROCS(1)

	wg := sync.WaitGroup{}

	nStart := thinkutils.DateTime.TimestampMs()
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go task(&wg)
	}

	wg.Wait()
	fmt.Printf("%d", thinkutils.DateTime.TimestampMs()-nStart)
}
