package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"fmt"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

type User struct {
	Name string
	Age  uint8
}

type AjaxResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []User `json:"data"`
}

func jsonObject() {
	user := User{Name: "aaa", Age: 100}

	log.Info(thinkutils.JSONUtils.ToJson(user))
}

func jsonArray() {
	lstUser := []User{{Name: "a", Age: 1}, {Name: "b", Age: 2}}
	log.Info(thinkutils.JSONUtils.ToJson(lstUser))
}

func fromjson() {
	szJson := `{"Name":"aaa","Age":100}`
	var user User
	err := thinkutils.JSONUtils.FromJson(szJson, &user)
	if err != nil {
		fmt.Println(err)
	}
}

func parseAjaxResult() {
	szJson := `{"code": 200, "msg":"success", "data": [{"Name":"a","Age":1},{"Name":"b","Age":2}]}`
	var ajaxRet AjaxResult
	err := thinkutils.JSONUtils.FromJson(szJson, &ajaxRet)
	if err != nil {
		fmt.Println(err)
	}
}

func fromjsonArray() {
	szJson := `[{"Name":"a","Age":1},{"Name":"b","Age":2}]`
	var user []User
	err := thinkutils.JSONUtils.FromJson(szJson, &user)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	jsonObject()
	jsonArray()

	fromjson()
	fromjsonArray()
	parseAjaxResult()
}
