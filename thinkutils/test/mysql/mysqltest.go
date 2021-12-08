package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"fmt"
	"sync"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func basicQueryJSON(wg *sync.WaitGroup) {
	db := thinkutils.ThinkMysql.QuickConn()

	rows, err := db.Query("SELECT * FROM sys_user")
	if err != nil {
		return
	}

	szRet := thinkutils.ThinkMysql.ToJSON(rows)
	fmt.Println(szRet)
	wg.Done()
	defer rows.Close()
}

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go basicQueryJSON(&wg)
	}

	wg.Wait()
	//time.Sleep(10 * time.Second)
	//basicQueryJSON()
}
