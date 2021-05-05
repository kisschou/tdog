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
		go NewLogger().Error(err.Error())
		return output
	}
	return output
}

func (excel *Excel) Open() (excelImpl *xlsx.File) {
	file := excel.Path + "/" + excel.File
	var err error
	excelImpl, err = xlsx.OpenFile(file)
	if err != nil {
		go NewLogger().Error(err.Error())
		return
	}
	return
}
