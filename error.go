package tdog

import (
	"path/filepath"
)

type (
	err struct {
		code string
		msg  string
	}
)

func NewError(input string) *err {
	errInfo := make([]string, 0)
	errFile := NewUtil().GetEnv("ERROR_FILE")
	configImpl := NewConfig().SetPath(filepath.Dir(errFile["ERROR_FILE"])).SetFile(filepath.Base(errFile["ERROR_FILE"]))
	if configImpl.Get(input).IsExists() {
		errInfo = configImpl.Get(input).ToStringSlice()
	} else {
		errInfo = configImpl.Get("BASE.UNKNOW").ToStringSlice()
	}
	return &err{code: errInfo[0], msg: errInfo[1]}
}

func (e *err) Code() string {
	return e.code
}

func (e *err) Msg() string {
	return e.msg
}
