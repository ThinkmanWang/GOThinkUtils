package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	thinkutils.FileUtils.ReadLine("/Users/wangxiaofeng/Github-Thinkman/GolandProjects/GOThinkUtils/test.txt", func(szLine string) {
		log.Info(szLine)
	})
}
