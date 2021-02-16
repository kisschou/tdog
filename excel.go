package tdog

import (
	"github.com/tealeg/xlsx"
)

type Excel struct {
	Path     string
	File     string
	SheetNum int
}

func (excel *Excel) Get() [][][]string {
	file := excel.Path + "/" + excel.File
	output, err := xlsx.FileToSlice(file)
	if err != nil {
		LogTdog := new(Logger)
		LogTdog.New(err.Error())
		return output
	}
	return output
}

func (excel *Excel) Open() (excelImpl *xlsx.File) {
	file := excel.Path + "/" + excel.File
	var err error
	excelImpl, err = xlsx.OpenFile(file)
	if err != nil {
		LogTdog := new(Logger)
		LogTdog.New(err.Error())
		return
	}
	return
}