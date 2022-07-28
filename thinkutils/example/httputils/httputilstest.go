package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"runtime"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func GetJSONTest()  {
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
	log.Info(thinkutils.JSONUtils.TrimJSON(szJson))
}

func TrimJSONTest()  {
	szJson := `
		{
			"maxCR": 0.65,\n\n\n
			"minCR": 0.55,
			"nativeSlotIds": ["575076", "575075", "575077"],
			"loadRnd": 30,
			"loadInv": 5,
			"rewardSlotId": "533213",
			"nonAuditVersion": [],
			"clickInv": 90,
			"closeJump": 0,
			"greedy": 0,
			"fatigue": true,
			"retentionOn": false,
			"bg": false
			"a": [
				{"a": 1}
			]
		}
	`
	log.Info(szJson)
	log.Info(thinkutils.JSONUtils.TrimJSON(szJson))
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    log.Info("Hello World")

    TrimJSONTest()
    GetJSONTest()
}
