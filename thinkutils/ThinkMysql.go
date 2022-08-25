package thinkutils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
	"strings"
	"sync"
	"time"
)

type thinkmysql struct {
}

var (
	g_lockMysql  sync.Mutex
	g_mapMysqlDB map[string]*sql.DB
)

func (this thinkmysql) makeConn(szHost string,
	nPort int,
	szUser string,
	szPwd string,
	szDb string,
	nMaxConn int) *sql.DB {
	defer g_lockMysql.Unlock()
	g_lockMysql.Lock()

	szConn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", szUser, szPwd, szHost, nPort, szDb)
	db := g_mapMysqlDB[szConn]
	if nil == db {
		_db, err := sql.Open("mysql", szConn)
		if err != nil {
			return nil
		}

		_db.SetConnMaxLifetime(time.Minute * 3)
		_db.SetMaxOpenConns(nMaxConn)
		_db.SetMaxIdleConns(2)

		g_mapMysqlDB[szConn] = _db
		db = _db
	}

	return db
}

func (this thinkmysql) Conn(szHost string,
	nPort int,
	szUser string,
	szPwd string,
	szDb string,
	nMaxConn int) *sql.DB {

	//id:password@tcp(your-amazonaws-uri.com:3306)/dbname
	szConn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", szUser, szPwd, szHost, nPort, szDb)
	//szConn := fmt.Sprintf("%s:%s@%s:%d/%s", szUser, szPwd, szHost, nPort, szDb)

	if nil == g_mapMysqlDB {
		g_mapMysqlDB = make(map[string]*sql.DB)
	}

	pDb := g_mapMysqlDB[szConn]
	if nil == pDb {
		pDb = this.makeConn(szHost, nPort, szUser, szPwd, szDb, nMaxConn)
	}

	//log.Info("%p %p", g_mapMysqlDB, pDb)
	return pDb
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

//func (this thinkmysql) ScanRow(rows *sql.Rows, dest interface{}) error {
//	s := reflect.ValueOf(dest).Elem()
//	numCols := s.NumField()
//	columns := make([]interface{}, numCols)
//	for i := 0; i < numCols; i++ {
//		field := s.Field(i)
//		columns[i] = field.Addr().Interface()
//	}
//
//	err := rows.Scan(columns...)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (this thinkmysql) ScanRow(rows *sql.Rows, dest interface{}) error {
	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}

	//s := reflect.ValueOf(dest).Elem()
	//numCols := s.NumField()
	columns := make([]interface{}, len(columnNames))
	for i := 0; i < len(columnNames); i++ {
		addr, bExists := StructUtils.FieldAddrByTag(dest, "field", columnNames[i])
		if bExists {
			columns[i] = addr
		} else {
			columns[i] = new(interface{})
		}
	}

	err = rows.Scan(columns...)
	if err != nil {
		return err
	}

	return nil
}

func (this thinkmysql) LastInsertId(tx *sql.Tx) (int64, error) {
	if nil == tx {
		return 0, errors.New("tx could not be null")
	}

	var nId int64

	if err := tx.QueryRow(`select LAST_INSERT_ID()`).Scan(&nId); err != nil {
		return 0, err
	}

	return nId, nil
}

func (this thinkmysql) TxBegin(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		return nil
	}

	return tx
}

func (this thinkmysql) TxRollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		log.Error(err.Error())
	}
}

func (this thinkmysql) TxCommit(tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

type NullBool struct {
	sql.NullBool
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if nb.Valid {
		return json.Marshal(nb.Bool)
	}
	return json.Marshal(nil)
}

func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}
	return nil
}

// Nullable Float64 that overrides sql.NullFloat64
type NullFloat64 struct {
	sql.NullFloat64
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(nil)
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}
	return nil
}

// Nullable Int64 that overrides sql.NullInt64
type NullInt64 struct {
	sql.NullInt64
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int64)
	}
	return json.Marshal(nil)
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

// Nullable String that overrides sql.NullString
type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}
