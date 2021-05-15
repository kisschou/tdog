package tdog

import (
	"github.com/tealeg/xlsx"
)

type excel struct {
	file     string
	sheetNum int
}

func NewExcel(file string) *excel {
	return &excel{file: file}
}

func (e *excel) Get() [][][]string {
	file := e.file
	output, err := xlsx.FileToSlice(file)
	if err != nil {
		go NewLogger().Error(err.Error())
		return output
	}
	return output
}

func (e *excel) Open() (excelImpl *xlsx.File) {
	file := e.file
	var err error
	excelImpl, err = xlsx.OpenFile(file)
	if err != nil {
		go NewLogger().Error(err.Error())
		return
	}
	return
}
