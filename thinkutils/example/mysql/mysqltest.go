package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
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
	Id thinkutils.NullInt64 `json:"id" field:"id"`
	Name thinkutils.NullString `json:"name" field:"name"`
}

func basicQueryJSON(wg *sync.WaitGroup) {

	db := thinkutils.ThinkMysql.QuickConn()
	defer wg.Done()

	rows, err := db.Query(`
		SELECT 
       		id, name 
		FROM 
		    t_test`)
	if err != nil {
		return
	}
	defer rows.Close()

	//lstRet := make([]MyType, 1)

	for rows.Next() {
		game := MyType{}

		err = thinkutils.ThinkMysql.ScanRow(rows, &game)
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.Info("%s", thinkutils.JSONUtils.ToJson(game))
	}
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
