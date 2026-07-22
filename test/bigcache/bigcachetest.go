package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	if err := thinkutils.ThinkBigCacheInstance().Start(); err != nil {
		log.Error("Failed to start big cache instance : %s", err.Error())
		return
	}

	if data, err := thinkutils.ThinkBigCacheInstance().Get("test001"); err != nil {
		log.Error("Failed to get test001 : %s", err.Error())
	} else {
		log.Info("Test001 : %s", string(data))
	}

	_ = thinkutils.ThinkBigCacheInstance().Set("test001", []byte("Hello World"))
	thinkutils.ThinkBigCacheInstance().AddUpdateListener(thinkutils.UPDATE_1_MIN, func() {
		log.Info("Update every 1 min")
		if data, err := thinkutils.ThinkBigCacheInstance().Get("test001"); err != nil {
			log.Error("Failed to get test001 : %s", err.Error())
		} else {
			log.Info("Test001 : %s", string(data))
		}
	})
	thinkutils.ThinkBigCacheInstance().AddUpdateListener(thinkutils.UPDATE_5_MIN, func() {
		log.Info("Update every 5 min")
		if data, err := thinkutils.ThinkBigCacheInstance().Get("test001"); err != nil {
			log.Error("Failed to get test001 : %s", err.Error())
		} else {
			log.Info("Test001 : %s", string(data))
		}
	})
	thinkutils.ThinkBigCacheInstance().AddUpdateListener(thinkutils.UPDATE_10_MIN, func() {
		log.Info("Update every 10 min")
		if data, err := thinkutils.ThinkBigCacheInstance().Get("test001"); err != nil {
			log.Error("Failed to get test001 : %s", err.Error())
		} else {
			log.Info("Test001 : %s", string(data))
		}
	})

	select {}
}
