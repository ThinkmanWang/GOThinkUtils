package thinkutils

type ThinkError struct {
	error
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (this *ThinkError) Error() string {
	return this.Msg
}

//func NewThinkError(format string, a ...any) *ThinkError {
//	szTxt := fmt.Sprintf(format, a...)
//	return &ThinkError{
//		Code: 500,
//		Msg:  szTxt,
//	}
//}
//
//func NewThinkErrorEx(nCode int64, format string, a ...any) *ThinkError {
//	szTxt := fmt.Sprintf(format, a...)
//	return &ThinkError{
//		Code: nCode,
//		Msg:  szTxt,
//	}
//}
