package thinkutils

import (
	"errors"
	"fmt"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/qqwry"
	"strings"
)

type qqwryutils struct {
}

func (this qqwryutils) downloadQqwry(szAppCode string) (string, error) {
	szResp, err := HttpUtils.GetWithHeader("http://cz88.rtbasia.com/getDownloadInfo", map[string]string{
		"Authorization": fmt.Sprintf("APPCODE %s", szAppCode),
	})
	if err != nil {
		//log.Error(err.Error())
		return "", err
	}

	log.Info(szResp)

	dictRet := make(map[string]any)
	if err = JSONUtils.FromJson(szResp, &dictRet); err != nil {
		//log.Error(err.Error())
		return "", err
	}

	nCode, bExist := dictRet["code"]
	if false == bExist {
		return "", errors.New("return code not exists")
	}

	if 200 != int32(nCode.(float64)) {
		return "", errors.New("ret code incorrect")
	}

	dictData, bExistsData := dictRet["data"]
	if false == bExistsData {
		return "", errors.New("data not in return json")
	}

	szDownloadUrl, bExistsUrl := dictData.(map[string]any)["downloadUrl"]
	if false == bExistsUrl {
		return "", errors.New("downloadUrl not in return JSON")
	}

	log.Info(szDownloadUrl.(string))

	szDateTime := DateTime.CurDatetime()
	szDateTime = strings.ReplaceAll(szDateTime, "-", "")
	szDateTime = strings.ReplaceAll(szDateTime, " ", "")
	szDateTime = strings.ReplaceAll(szDateTime, ":", "")

	szFileName := fmt.Sprintf("%s.dat", szDateTime)
	err = HttpUtils.DownloadFile(szFileName, szDownloadUrl.(string))
	if err != nil {
		return "", err
	}

	return szFileName, nil
}

func (this qqwryutils) Init(szAppCode string) error {
	szFile, err := this.downloadQqwry(szAppCode)
	if err != nil {
		return err
	}

	pDat := qqwry.NewQQwry(szFile)
	GetMemCacheInstance().Set("QQWRY_DAT", 120, pDat, szAppCode, this.RefreshData)

	return nil
}

func (this qqwryutils) RefreshData(pUserData any) error {
	if nil == pUserData {
		return nil
	}

	szAppCode := pUserData.(string)
	szFile, err := this.downloadQqwry(szAppCode)
	if err != nil {
		log.Error("**QQWRY** %s", err.Error())
		return err
	}

	pDat := qqwry.NewQQwry(szFile)
	GetMemCacheInstance().Set("QQWRY_DAT", 120, pDat, szAppCode, this.RefreshData)

	return nil
}

func (this qqwryutils) IPLocation(szIP string) (*qqwry.Rq, error) {
	pData := GetMemCacheInstance().Get("QQWRY_DAT")
	if nil == pData {
		return nil, errors.New("QQWRY not found")
	}

	pDat := pData.(*qqwry.QQwry)
	return pDat.Find(szIP), nil
}
