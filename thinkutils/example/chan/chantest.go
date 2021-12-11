package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func worker(ch chan string) {
	for {
		szTxt := <-ch
		log.Info(szTxt)
	}
}

func main() {
	ch := make(chan string)
	go worker(ch)

	for {
		ch <- thinkutils.DateTime.CurDatetime()
		time.Sleep(time.Second)
	}
}
