package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	if err := thinkutils.QQwryUtils.Init("f4672ec4a55a436c87f7ce1f1a631ea9"); err != nil {
		log.Error(err.Error())
		return
	}

	if err := thinkutils.QQwryUtils.RefreshData("f4672ec4a55a436c87f7ce1f1a631ea9"); err != nil {
		log.Error(err.Error())
		//return
	}

	pRet, err := thinkutils.QQwryUtils.IPLocation("182.139.183.98")
	if err != nil {
		log.Error(err.Error())
		return
	}

	if nil == pRet {
		log.Error("parse ip error, return null")
		return
	}

	log.Info("%s", thinkutils.JSONUtils.ToJson(pRet))
}
