package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"runtime"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    log.Info("Hello World")

    szJson, err := thinkutils.HttpUtils.Get("https://czy7n.jiegames.com/jy-game/qmfcwAPK/OppoApkGameConfig.json")
    if err != nil {
    	log.Info(err.Error())
		return
	}

	if thinkutils.StringUtils.IsEmpty(szJson) {
		log.Error("ret empty")
		return
	}

	//szJson = strings.Replace(szJson, "\n", "", -1)
	//szJson = strings.Replace(szJson, "\t", "", -1)
	bIsJson := thinkutils.JSONUtils.IsJSONString(szJson)
	if false == bIsJson {
		log.Info("ret value is not json")
		return
	}
	log.Info(szJson)
}
