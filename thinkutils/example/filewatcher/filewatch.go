package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func onCreate(szFile string) {
	log.Info(szFile)
}

func onModify(szFile string) {
	log.Info(szFile)
}

func onDelete(szFile string) {
	log.Info(szFile)
}

func main() {
	pWatch := &thinkutils.ThinkNotify{
		OnCreate: onCreate,
		OnModify: onModify,
		OnDelete: onDelete,
	}

	//log.Info("%p", pWatch)
	pWatch.Watch("/Users/wangxiaofeng/Github-Thinkman/GolandProjects/github.com/ThinkmanWang/GOThinkUtils/temp")
}
