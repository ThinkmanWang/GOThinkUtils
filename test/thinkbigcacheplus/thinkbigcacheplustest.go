package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func commonStrToByte(val any) ([]byte, error) {
	return []byte(val.(string)), nil
}

func commonByteToStr(data []byte) (any, error) {
	return string(data), nil
}

func testSave() {
	thinkutils.ThinkBigCachePlusInstance().RegByteFunction("commonStr", commonStrToByte, commonByteToStr)
	_ = thinkutils.ThinkBigCachePlusInstance().Start()

	_ = thinkutils.ThinkBigCachePlusInstance().Set("commonStr", "hello", "world123")
	_ = thinkutils.ThinkBigCachePlusInstance().SaveToDisk()
}

func testLoad() {
	thinkutils.ThinkBigCachePlusInstance().RegByteFunction("commonStr", commonStrToByte, commonByteToStr)
	_ = thinkutils.ThinkBigCachePlusInstance().Start()

	if data, err := thinkutils.ThinkBigCachePlusInstance().Get("commonStr", "hello"); err != nil {
		log.Error("Failed to get hello from cache : %s", err.Error())
	} else {
		log.Info("%s", data.(string))
	}
}

func main() {
	//testSave()
	testLoad()
}
