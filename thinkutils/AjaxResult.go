package thinkutils

type AjaxResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func AjaxResultSuccess() *AjaxResult {
	return &AjaxResult{Code: 0, Msg: "success"}
}

func AjaxResultSuccessWithData(data interface{}) *AjaxResult {
	return &AjaxResult{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}

func AjaxResultError() *AjaxResult {
	return &AjaxResult{Code: 500, Msg: "Server Error"}
}

func AjaxResultErrorWithMsg(msg string) *AjaxResult {
	return &AjaxResult{
		Code: 500,
		Msg: msg,
	}
}
