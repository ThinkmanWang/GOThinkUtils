package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"fmt"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

type User struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

type AjaxResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []User `json:"data"`
}

func jsonObject() {
	user := User{Name: "aaa", Age: 100}

	log.Info(thinkutils.JSONUtils.ToJson(&user))
}

func jsonArray() {
	lstUser := []User{{Name: "a", Age: 1}, {Name: "b", Age: 2}}
	log.Info(thinkutils.JSONUtils.ToJson(lstUser))
}

func fromjson() {
	szJson := `{"name":"aaa","age":100}`
	var user User
	err := thinkutils.JSONUtils.FromJson(szJson, &user)
	if err != nil {
		fmt.Println(err)
	}
}

func parseAjaxResult() {
	szJson := `{"code": 200, "msg":"success", "data": [{"name":"a","age":1},{"name":"b","age":2}]}`
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

func testAjaxResult() {
	pAjax := thinkutils.AjaxResultSuccess()

	szJson := `[{"Name":"a","Age":1},{"Name":"b","Age":2}]`
	var user []User
	err := thinkutils.JSONUtils.FromJson(szJson, &user)
	if err != nil {
		fmt.Println(err)
	}
	pAjax.Data = user

	log.Info(thinkutils.JSONUtils.ToJson(pAjax))
}

func testAjaxResult1() {

	szJson := `[{"Name":"a","Age":1},{"Name":"b","Age":2}]`
	var user []User
	err := thinkutils.JSONUtils.FromJson(szJson, &user)
	if err != nil {
		fmt.Println(err)
	}

	pAjax := thinkutils.AjaxResultSuccessWithData(user)

	log.Info(thinkutils.JSONUtils.ToJson(pAjax))
}

func testAjaxResult2() {
	log.Info(thinkutils.JSONUtils.ToJson(thinkutils.AjaxResultSuccess()))
}

func main() {
	jsonObject()
	jsonArray()

	fromjson()
	fromjsonArray()
	parseAjaxResult()

	testAjaxResult()
	testAjaxResult1()
	testAjaxResult2()
}
