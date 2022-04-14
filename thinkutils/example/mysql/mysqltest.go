package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
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

type MyType struct {
	AppId uint64 `json:"id" field:"app_id"`
	Name string `json:"name" field:"name"`
}

func basicQueryJSON(wg *sync.WaitGroup) {

	db := thinkutils.ThinkMysql.QuickConn()

	rows, err := db.Query(`
		SELECT 
       		* 
		FROM 
		    t_game`)
	if err != nil {
		return
	}

	//lstRet := make([]MyType, 1)

	for rows.Next() {
		game := MyType{}

		err = thinkutils.ThinkMysql.ScanRow(rows, &game)
		if err != nil {
			return
		}

		fmt.Println(game)
	}


	defer rows.Close()
	defer wg.Done()

}

func main() {
	wg := sync.WaitGroup{}
	//for i := 0; i < 100; i++ {
	wg.Add(1)
	go basicQueryJSON(&wg)
	//}

	wg.Wait()
	//time.Sleep(10 * time.Second)
	//basicQueryJSON()
}
