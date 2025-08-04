package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	szText := "Hello World"
	szOut, _ := thinkutils.AESUtils.EncryptString(szText, "yxt2hibpi5dlbi8y")
	log.Info(szOut)

	szOut, _ = thinkutils.AESUtils.DecryptString(szOut, "yxt2hibpi5dlbi8y")
	log.Info(szOut)
}
