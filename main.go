package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"fmt"
	"github.com/buger/jsonparser"
	"strconv"
	"time"
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

func cor1() {
	time.Sleep(1000)
	fmt.Println("FXXK")
}

func cor2(chRet chan string) {
	chRet <- "FXXXXXXXXXXXXXK"
}

func coTest() {
	go cor1()

	c := make(chan string)
	go cor2(c)

	szRet := <-c
	fmt.Println(szRet)
}

func logTest() {
	//logger.SetLogger(`{"Console": {"level": "DEBG"}`)
	log.Info("FXXK")
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

func jsonTest() {
	data := []byte(`{
	  "person": {
		"name": {
		  "first": "Leonid",
		  "last": "Bugaev",
		  "fullName": "Leonid Bugaev"
		},
		"github": {
		  "handle": "buger",
		  "followers": 109
		},
		"avatars": [
		  { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
		]
	  },
	  "company": {
		"name": "Acme"
	  }
	}`)

	szTxt, _ := jsonparser.GetString(data, "person", "name", "fullName")
	log.Info(szTxt)

	nFollower, _ := jsonparser.GetInt(data, "person", "github", "followers")
	log.Info("%d", nFollower)

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		szUrl, _ := jsonparser.GetString(value, "url")
		log.Info(szUrl)
	}, "person", "avatars")
}

func main() {
	log.Info("%d", 123)
	fmt.Println("Hello World")

	//var logger *log.Logger = new(log.Logger)

	datetimeTest()

	coTest()
	fmt.Println("fxxk1")

	logTest()

	log.Info(thinkutils.RandUtils.RandPasssword(8))
	log.Info(thinkutils.RandUtils.UUID())

	md5Test()

	log.Info(thinkutils.IPUtils.LocalIP())

	jsonTest()

	time.Sleep(10 * 1000)
}
