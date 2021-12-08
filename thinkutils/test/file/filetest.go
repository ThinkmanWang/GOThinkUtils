package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	thinkutils.FileUtils.ReadLine("test.txt", func(nLine uint32, szLine string) {
		log.Info("%d %s", nLine, szLine)
	})
}
