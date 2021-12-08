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

/*
	sqlStatement := `
	UPDATE users
	SET first_name = $2, last_name = $3
	WHERE id = $1;`
	res, err := db.Exec(sqlStatement, 5, "NewFirst", "NewLast")
	if err != nil {
	  panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
	  panic(err)
	}
	fmt.Println(count)
*/
func basicQueryJSON(wg *sync.WaitGroup) {
	db := thinkutils.ThinkMysql.QuickConn()

	rows, err := db.Query(`
		SELECT 
       		* 
		FROM 
		    sys_user 
		WHERE 
		    user_id > ?`, 0)
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
