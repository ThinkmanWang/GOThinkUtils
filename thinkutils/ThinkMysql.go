package thinkutils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
	"strings"
	"sync"
	"time"
)

type thinkmysql struct {
	m_lock  sync.Mutex
	m_mapDB map[string]*sql.DB
}

func (this thinkmysql) Conn(szHost string,
	nPort int,
	szUser string,
	szPwd string,
	szDb string,
	nMaxConn int) *sql.DB {
	defer this.m_lock.Unlock()
	this.m_lock.Lock()

	//id:password@tcp(your-amazonaws-uri.com:3306)/dbname
	szConn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", szUser, szPwd, szHost, nPort, szDb)
	//szConn := fmt.Sprintf("%s:%s@%s:%d/%s", szUser, szPwd, szHost, nPort, szDb)

	if nil == this.m_mapDB {
		this.m_mapDB = make(map[string]*sql.DB)
	}

	pDb := this.m_mapDB[szConn]
	if pDb != nil {
		return pDb
	}

	db, err := sql.Open("mysql", szConn)
	if err != nil {
		return nil
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(nMaxConn)
	db.SetMaxIdleConns(2)

	this.m_mapDB[szConn] = db

	return db
}

func (this thinkmysql) QuickConn() *sql.DB {
	cfg, err := ini.Load("app.ini")
	if err != nil {
		return this.Conn("127.0.0.1", 3306, "root", "123456", "db1", 16)
	}

	return this.Conn(cfg.Section("mysql").Key("host").String(),
		cfg.Section("mysql").Key("port").MustInt(),
		cfg.Section("mysql").Key("user").String(),
		cfg.Section("mysql").Key("password").String(),
		cfg.Section("mysql").Key("db").String(),
		cfg.Section("mysql").Key("max_conn").MustInt())
}

func (this thinkmysql) ToJSON(rows *sql.Rows) string {
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return ""
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {

		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {

			switch strings.ToUpper(v.DatabaseTypeName()) {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT4", "INT", "BIGINT", "INTEGER", "TINYINT":
				scanArgs[i] = new(sql.NullInt64)
				break
			case "DOUBLE", "FLOAT", "DECIMAL":
				scanArgs[i] = new(sql.NullFloat64)
				break
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)

		if err != nil {
			return ""
		}

		masterData := map[string]interface{}{}

		for i, v := range columnTypes {

			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}

	z, err := json.Marshal(finalRows)

	szJson := StringUtils.BytesToString(z)

	return szJson
}
