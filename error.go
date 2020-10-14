package tdog

import (
)

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
