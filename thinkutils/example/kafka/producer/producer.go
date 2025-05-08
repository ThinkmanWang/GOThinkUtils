package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {

	for i := 0; i < 10000; i++ {
		go func() {
			//szMsg := fmt.Sprintf("[%d] %s", nIndex, thinkutils.DateTime.CurDatetime())
			dictDetail := map[string]any{
				"uid":        "00001",
				"logtype":    "loadAd",
				"adPosition": "aaaaabbbsss",
				"timestamp":  thinkutils.DateTime.TimestampMs(),
			}
			dictMsg := map[string]any{
				"json": dictDetail,
			}

			thinkutils.KafkaUtils.SendMsg("172.16.0.106:9092,172.16.8.106:9092,172.16.12.54:9092",
				"topic-game_loghub_common",
				[]byte(thinkutils.JSONUtils.ToJson(dictMsg)))

			log.Info("Send %s", thinkutils.JSONUtils.ToJson(dictMsg))
			time.Sleep(time.Duration(500) * time.Millisecond)
		}()

	}

	time.Sleep(60 * time.Second)
}
