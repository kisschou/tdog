package tdog

import (
	"reflect"
	"testing"

	"github.com/tealeg/xlsx"
)

func TestExcel_Get(t *testing.T) {
	type fields struct {
		Path     string
		File     string
		SheetNum int
	}

	tests := []struct {
		name   string
		fields fields
		want   [][][]string
	}{
		// TODO: Add test cases.
		{
			name: "测试excel读取",
			fields: fields{
				Path:     "/Users/kisschou/data/golang/src/all-service/sources/FileService/data/excel",
				File:     "小区信息表.xlsx",
				SheetNum: 0,
			},
			want: [][][]string{
				{
					{
						"物业项目名称(即小区名称)",
						"行政区名称",
						"物业项目地址",
						"自然幢数",
						"开发日期",
						"占地面积",
						"建筑面积",
					},
					{
						"康嘉御景A区",
						"城厢区",
						"城厢区龙桥街道溪南路388号",
						"36",
						"2011-03-28",
						"73790.27",
						"54721",
					},
					{
						"奥元雅居（PS拍-2013-02号）",
						"荔城区",
						"荔城区镇海街道丰美422弄",
						"3",
						"2015-03-27",
						"3286.7",
						"8848.41",
					},
					{
						"秀屿区配建经济适用房（A4、A5幢）",
						"秀屿区",
						"秀屿区笏石镇四新小区",
						"2",
						"2011-1-6",
						"19372",
						"18720.08",
					},
					{
						"中港实业锦江大楼",
						"荔城区",
						"福建省莆田市荔城区黄石镇万好街799号",
						"2",
						"2011-3-20",
						"6779.02",
						"130565.99",
					},
					{
						"莆田喜盈门建材家具广场",
						"城厢区",
						"莆田市城厢区荔园路北侧",
						"2",
						"-",
						"-",
						"165488.69",
					},
					{
						"汉庭花园B区",
						"荔城区",
						"荔城区汉庭路298号",
						"15",
						"-",
						"-",
						"-",
					},
					{
						"中海天下",
						"荔城区",
						"莆田市荔城区拱辰街道东圳东路1199号中海天下",
						"7",
						"2012-03-21",
						"37435.09",
						"198803.4",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			excel := &Excel{
				Path:     tt.fields.Path,
				File:     tt.fields.File,
				SheetNum: tt.fields.SheetNum,
			}
			if got := excel.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Excel.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExcel_Open(t *testing.T) {
	type fields struct {
		Path     string
		File     string
		SheetNum int
	}
	tests := []struct {
		name          string
		fields        fields
		wantExcelImpl *xlsx.File
	}{
		// TODO: Add test cases.
		{
			name: "测试excel打开",
			fields: fields{
				Path:     "/Users/kisschou/data/golang/src/all-service/sources/FileService/data/excel",
				File:     "小区信息表.xlsx",
				SheetNum: 0,
			},
			wantExcelImpl: new(xlsx.File),
		},
	}
	type estate struct {
		Name         string `xlsx:"0"`
		Region       string `xlsx:"1"`
		Address      string `xlsx:"2"`
		BuildingNum  string `xlsx:"3"`
		BuildingTime string `xlsx:"4"`
		FloorArea    string `xlsx:"5"`
		CoveredArea  string `xlsx:"6"`
	}
	estateList := make([]estate, 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			excel := &Excel{
				Path:     tt.fields.Path,
				File:     tt.fields.File,
				SheetNum: tt.fields.SheetNum,
			}
			if gotExcelImpl := excel.Open(); gotExcelImpl != nil {
				for i := 0; i < gotExcelImpl.Sheet["Sheet1"].MaxRow; i++ {
					estateInfo := new(estate)
					err := gotExcelImpl.Sheet["Sheet1"].Rows[i].ReadStruct(estateInfo)
					if err != nil {
						LogTdog := new(Logger)
						LogTdog.New(err.Error())
						continue
					}
					estateList = append(estateList, *estateInfo)
				}
				if len(estateList) != 8 {
					t.Errorf("Data length is fail")
				}
			} else {
				t.Errorf("gotExcelImpl is nil")
			}
		})
	}
}
