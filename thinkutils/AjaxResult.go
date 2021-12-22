package thinkutils

type AjaxResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func AjaxResultSuccess() AjaxResult {
	return AjaxResult{Code: 200, Msg: "success"}
}

func AjaxResultSuccessWithData(data interface{}) AjaxResult {
	return AjaxResult{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func AjaxResultError() AjaxResult {
	return AjaxResult{Code: 500, Msg: "Server Error"}
}
