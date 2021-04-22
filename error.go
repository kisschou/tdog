package tdog

import ()

type (
	Error struct {
	}
)

func (e *Error) GetError(errCode string) (errMsg string) {
	errMsg = errCode
	if NewConfig().Get("error." + errCode).IsExists() {
		errMsg = NewConfig().Get("error." + errCode).ToString()
	} else {
		errMsg = NewConfig().Get("error.ERROR_UNKNOW").ToString()
	}
	return
}

func (e *Error) GetErrorCode(errCode string) (code int) {
	if NewConfig().Get("error_map." + errCode).IsExists() {
		code = NewConfig().Get("error_map." + errCode).ToInt()
	} else {
		code = NewConfig().Get("error_map.ERROR_UNKNOW").ToInt()
	}
	return
}
