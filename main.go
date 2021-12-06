package main

import (
	"fmt"

	"GOThinkUtils/thinkutils"
)

func main() {
	fmt.Println("Hello World")

	fmt.Println(thinkutils.DateTime.Timestamp())
	fmt.Println(thinkutils.DateTime.TimestampMs())
	fmt.Println(thinkutils.StringUtils.IsEmpty(" 123 "))

	var pszTxt *string = new(string)
	*pszTxt = " 123"
	fmt.Println(thinkutils.StringUtils.IsEmptyPtr(pszTxt))
	fmt.Println(thinkutils.StringUtils.IsEmpty(*pszTxt))
}
