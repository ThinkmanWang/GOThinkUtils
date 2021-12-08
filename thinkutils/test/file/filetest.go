package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	thinkutils.FileUtils.ReadLine("/Users/wangxiaofeng/Github-Thinkman/GolandProjects/GOThinkUtils/test.txt", func(nLine uint32, szLine string) {
		log.Info("%d %s", nLine, szLine)
	})
}
