package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func RefreshData(pData any) error {
	log.Info("FXXK")
	thinkutils.GetMemCacheInstance().Set("fxxk", 120, nil, nil, RefreshData)
	return nil
}

func main() {
	log.Info("Hello World")

	thinkutils.GetMemCacheInstance().Start()
	thinkutils.GetMemCacheInstance().Set("fxxk", 120, nil, nil, RefreshData)

	for {
		time.Sleep(10 * time.Second)
	}
}
