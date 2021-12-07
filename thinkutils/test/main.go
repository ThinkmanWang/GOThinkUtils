package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"fmt"
	"strconv"
)

var log *logger.LocalLogger = logger.DefaultLogger()

func datetimeTest() {
	fmt.Println(thinkutils.DateTime.Timestamp())
	fmt.Println(thinkutils.DateTime.TimestampMs())
	fmt.Println(thinkutils.StringUtils.IsEmpty(" 123 "))

	var pszTxt *string = new(string)
	*pszTxt = " 123"
	fmt.Println(thinkutils.StringUtils.IsEmptyPtr(pszTxt))
	fmt.Println(thinkutils.StringUtils.IsEmpty(*pszTxt))

	fmt.Println(thinkutils.DateTime.CurDatetime())
	fmt.Println(thinkutils.DateTime.Yesterday())
	fmt.Println(thinkutils.DateTime.Tomorrow())
	fmt.Println(thinkutils.DateTime.TimeStampToDateTime(thinkutils.DateTime.Timestamp()))
	fmt.Println(thinkutils.DateTime.Hour())
	fmt.Println(strconv.Atoi("05"))

	fmt.Println(thinkutils.DateTime.DiffDate(-3))
	fmt.Println(thinkutils.DateTime.DiffDate(4))
	fmt.Println(thinkutils.DateTime.DateToTimestamp("2021-12-06"))
	fmt.Println(thinkutils.DateTime.FirstDayOfMonth("2021-10-20"))
	fmt.Println(thinkutils.DateTime.LastDayOfMonth("2021-03-01"))

	lstDate := thinkutils.DateTime.DateBetweenStartEnd("2021-12-01", "2021-12-10")
	for i := 0; i < len(lstDate); i++ {
		fmt.Println(lstDate[i])
	}

	for _, szDate := range lstDate {
		fmt.Println(szDate)
	}

	fmt.Println(thinkutils.DateTime.StartEndOfWeek("2021-12-16"))
}

func md5Test() {
	log.Info(thinkutils.MD5Utils.MD5String("HHH"))

	szMd5 := thinkutils.MD5Utils.MD5File("/Users/wangxiaofeng/Github-Thinkman/GolandProjects/GOThinkUtils/GOThinkUtils")
	log.Info(szMd5)

	chRet := make(chan string)
	go thinkutils.MD5Utils.MD5FileCor("/Users/wangxiaofeng/Github-Thinkman/GolandProjects/GOThinkUtils/GOThinkUtils", chRet)
	szMd5 = <-chRet
	log.Info(szMd5)
}

func main() {
	datetimeTest()

	log.Info(thinkutils.RandUtils.RandPasssword(8))
	log.Info(thinkutils.RandUtils.UUID())

	md5Test()

	log.Info(thinkutils.IPUtils.LocalIP())
}
