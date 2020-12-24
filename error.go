package tdog

import ()

type (
	Error struct {
	}
)

func (e *Error) GetError(errCode string) (errMsg string) {
	ConfigLib := new(Config)
	errMsg = errCode
	if ConfigLib.Get("error." + errCode).IsExists() {
		errMsg = ConfigLib.Get("error." + errCode).String()
	}
	return
}

func (e *Error) GetErrorCode(errCode string) (code int) {
	ConfigLib := new(Config)
	if ConfigLib.Get("error_map." + errCode).IsExists() {
		code = ConfigLib.Get("error_map." + errCode).Int()
	}
	return
}
