package main

import (
	"fmt"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func main() {
	szResp, err := thinkutils.HttpUtils.GetWithHeader("http://cz88.rtbasia.com/getDownloadInfo", map[string]string{
		"Authorization": fmt.Sprintf("APPCODE %s", "f4672ec4a55a436c87f7ce1f1a631ea9"),
	})
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info(szResp)

	dictRet := make(map[string]any)
	if err = thinkutils.JSONUtils.FromJson(szResp, &dictRet); err != nil {
		log.Error(err.Error())
		return
	}

	nCode, bExist := dictRet["code"]
	if false == bExist {
		return
	}

	if 200 != int32(nCode.(float64)) {
		log.Info("ret code incorrect")
		return
	}

	dictData, bExistsData := dictRet["data"]
	if false == bExistsData {
		return
	}

	szDownloadUrl, bExistsUrl := dictData.(map[string]any)["downloadUrl"]
	if false == bExistsUrl {
		return
	}

	log.Info(szDownloadUrl.(string))

	_ = thinkutils.HttpUtils.DownloadFile("qqwry.dat", szDownloadUrl.(string))
}
