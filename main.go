package main

import (
	"fmt"

	thinkutils "GOThinkUtils/gothinkutils"
)

func main() {
	fmt.Println("Hello World")

	fmt.Println(thinkutils.Datetime{}.Timestamp())
	fmt.Println(thinkutils.Datetime{}.TimestampMs())
}
