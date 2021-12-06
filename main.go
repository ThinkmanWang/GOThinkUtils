package main

import (
	"GOThinkUtils/thinkutils"
	"fmt"
	"log"
	"strconv"
)

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
}

func main() {
	fmt.Println("Hello World")

	//var logger *log.Logger = new(log.Logger)
	log.Println("FXXK")

	datetimeTest()
}
