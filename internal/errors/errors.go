package errors

import "encoding/json"

type APIError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

const (
	TokenInvalid  = 40000
	ParamsInvalid = 50204
	Unavilable    = 50205
)

var ErrMsgMap = map[int]string{
	ParamsInvalid: "参数缺失",
	TokenInvalid:  "非法请求",
	Unavilable:    "网络异常，请稍后重试",
}

func (e APIError) Error() string {
	b, _ := json.Marshal(e.Msg)
	return string(b)
}

func (e APIError) Data() APIError {
	return APIError{
		Code: e.Code,
		Msg:  e.Msg,
	}
}

func NewError(code int) APIError {
	return APIError{
		Code: code,
		Msg:  ErrMsgMap[code],
	}
}

func NewErrorWithMsg(code int, msg string) APIError {
	return APIError{
		Code: code,
		Msg:  msg,
	}
}
