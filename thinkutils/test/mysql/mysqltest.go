package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"fmt"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func basicQueryJSON() {
	db := thinkutils.ThinkMysql.QuickConn()

	rows, err := db.Query("SELECT * FROM sys_user")
	if err != nil {
		return
	}

	szRet := thinkutils.ThinkMysql.ToJSON(rows)
	fmt.Println(szRet)
	defer rows.Close()
}

func main() {
	basicQueryJSON()
}
