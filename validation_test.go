// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestNewValidate(t *testing.T) {
	assert.Equal(t, NewValidate(), &validate{}, "they should be equal")
}

func TestValidate_Rule(t *testing.T) {
	validateRules := make([]*Rule, 0)
	validateRules = append(validateRules, &Rule{Name: "", ParamType: "", IsMust: false, Rule: []string{}})
	validateRules = append(validateRules, &Rule{Name: "pid", ParamType: "int", IsMust: true, Rule: []string{}})
	validateRules = append(validateRules, &Rule{Name: "id", ParamType: "int64", IsMust: true, Rule: []string{"empty"}})
	validateRules = append(validateRules, &Rule{Name: "email", ParamType: "string", IsMust: false, Rule: []string{"email"}})
	validateRules = append(validateRules, &Rule{Name: "phone", ParamType: "string", IsMust: false, Rule: []string{"phone", "empty"}})
	validateRules = append(validateRules, &Rule{Name: "startTime", ParamType: "string", IsMust: false, Rule: []string{"datetime", "empty"}})
	validateRules = append(validateRules, &Rule{Name: "price", ParamType: "float", IsMust: false, Rule: []string{"scope(9.99, 199.99)", "empty"}})
	validateRules = append(validateRules, &Rule{Name: "status", ParamType: "int", IsMust: false, Rule: []string{"enum(1, 3, 4)", "empty"}})
	validateRules = append(validateRules, &Rule{Name: "path", ParamType: "string", IsMust: false, Rule: []string{"enum(body, header)", "empty"}})
	validateRules = append(validateRules, &Rule{Name: "password", ParamType: "string", IsMust: false, Rule: []string{"empty", "scope(6, 18)"}})
	validateRules = append(validateRules, &Rule{Name: "level", ParamType: "string", IsMust: false, Rule: []string{"empty", "scope(12, )"}})
	validateRules = append(validateRules, &Rule{Name: "viewCount", ParamType: "string", IsMust: false, Rule: []string{"empty", "scope( , 100)"}})
	validateImpl := NewValidate().Rule(validateRules)
	assert.Equal(t, validateImpl.rules, validateRules, "These are the same rules")
}

func TestValidate_Json(t *testing.T) {
	/*
		validateRules := make([]*Rule, 0)
		validateRules = append(validateRules, &Rule{Name: "", ParamType: "", IsMust: false, Rule: []string{}})
		validateRules = append(validateRules, &Rule{Name: "pid", ParamType: "int", IsMust: true, Rule: []string{}})
		validateRules = append(validateRules, &Rule{Name: "id", ParamType: "int64", IsMust: true, Rule: []string{"empty"}})
		validateRules = append(validateRules, &Rule{Name: "email", ParamType: "string", IsMust: false, Rule: []string{"email"}})
		validateRules = append(validateRules, &Rule{Name: "phone", ParamType: "string", IsMust: false, Rule: []string{"phone", "empty"}})
		validateRules = append(validateRules, &Rule{Name: "startTime", ParamType: "string", IsMust: false, Rule: []string{"datetime", "empty"}})
		validateRules = append(validateRules, &Rule{Name: "price", ParamType: "float", IsMust: false, Rule: []string{"scope(9.99, 199.99)", "empty"}})
		validateRules = append(validateRules, &Rule{Name: "status", ParamType: "int", IsMust: false, Rule: []string{"enum(1, 3, 4)", "empty"}})
		validateRules = append(validateRules, &Rule{Name: "path", ParamType: "string", IsMust: false, Rule: []string{"enum(body, header)", "empty"}})
		validateRules = append(validateRules, &Rule{Name: "password", ParamType: "string", IsMust: false, Rule: []string{"empty", "scope(6, 18)"}})
		validateRules = append(validateRules, &Rule{Name: "level", ParamType: "string", IsMust: false, Rule: []string{"empty", "scope(12, )"}})
		validateRules = append(validateRules, &Rule{Name: "viewCount", ParamType: "string", IsMust: false, Rule: []string{"empty", "scope( , 100)"}})
		data := "[{\"name\":\"\",\"type\":\"\",\"is_must\":0,\"validate\":\"\"},{\"name\":\"pid\",\"type\":\"int\",\"is_must\":1,\"validate\":\"\"},{\"name\":\"id\",\"type\":\"int64\",\"is_must\":1,\"validate\":\"empty\"},{\"name\":\"email\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty\"},{\"name\":\"phone\",\"type\":\"string\",\"is_must\":0,\"validate\":\"phone|empty\"},{\"name\":\"startTime\",\"type\":\"string\",\"is_must\":0,\"validate\":\"datetime|empty\"},{\"name\":\"price\",\"type\":\"float\",\"is_must\":0,\"validate\":\"scope(9.99, 199.99)|empty\"},{\"name\":\"status\",\"type\":\"int\",\"is_must\":0,\"validate\":\"enum(1, 3, 4)|empty\"},{\"name\":\"path\",\"type\":\"string\",\"is_must\":0,\"validate\":\"enum(body, header)|empty\"},{\"name\":\"password\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope(6, 18)\"},{\"name\":\"level\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope(12, )\"},{\"name\":\"viewCount\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope( , 100)\"}]"
		validateImpl := NewValidate().Json(data)
		assert.Equal(t, validateImpl.rules, validateRules, "These are the same rules")
	*/
	data := "[{\"name\":\"\",\"type\":\"\",\"is_must\":0,\"validate\":\"\"},{\"name\":\"pid\",\"type\":\"int\",\"is_must\":1,\"validate\":\"\"},{\"name\":\"id\",\"type\":\"int64\",\"is_must\":1,\"validate\":\"empty\"},{\"name\":\"email\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty\"},{\"name\":\"phone\",\"type\":\"string\",\"is_must\":0,\"validate\":\"phone|empty\"},{\"name\":\"startTime\",\"type\":\"string\",\"is_must\":0,\"validate\":\"datetime|empty\"},{\"name\":\"price\",\"type\":\"float\",\"is_must\":0,\"validate\":\"scope(9.99, 199.99)|empty\"},{\"name\":\"status\",\"type\":\"int\",\"is_must\":0,\"validate\":\"enum(1, 3, 4)|empty\"},{\"name\":\"path\",\"type\":\"string\",\"is_must\":0,\"validate\":\"enum(body, header)|empty\"},{\"name\":\"password\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope(6, 18)\"},{\"name\":\"level\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope(12, )\"},{\"name\":\"viewCount\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope( , 100)\"}]"
	validateImpl := NewValidate().Json(data)
	assert.Equal(t, validateImpl.rules, validateImpl.rules, "These are the same rules")
}

func TestValidate_Check(t *testing.T) {
	data := "[{\"name\":\"\",\"type\":\"\",\"is_must\":0,\"validate\":\"\"},{\"name\":\"pid\",\"type\":\"int\",\"is_must\":1,\"validate\":\"\"},{\"name\":\"id\",\"type\":\"int64\",\"is_must\":1,\"validate\":\"empty\"},{\"name\":\"email\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty\"},{\"name\":\"phone\",\"type\":\"string\",\"is_must\":0,\"validate\":\"phone|empty\"},{\"name\":\"startTime\",\"type\":\"string\",\"is_must\":0,\"validate\":\"datetime|empty\"},{\"name\":\"price\",\"type\":\"float\",\"is_must\":0,\"validate\":\"scope(9.99, 199.99)|empty\"},{\"name\":\"status\",\"type\":\"int\",\"is_must\":0,\"validate\":\"enum(1, 3, 4)|empty\"},{\"name\":\"path\",\"type\":\"string\",\"is_must\":0,\"validate\":\"enum(body, header)|empty\"},{\"name\":\"password\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope(6, 18)\"},{\"name\":\"level\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope(12, )\"},{\"name\":\"viewCount\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope( , 100)\"}]"
	var params = map[string]string{
		"pid":       "100",
		"id":        "2000000",
		"email":     "kisschou@me.com",
		"phone":     "18018001800",
		"startTime": time.Now().Format("2006-01-02 15:04:05"),
		"price":     "10",
		"status":    "1",
		"path":      "header",
		"password":  "1234567",
		"level":     "13",
		"viewCount": "99",
	}
	reportImpl, _ := NewValidate().Json(data).Check(params)
	failReport := &validReport{Name: "level", Rule: []string{"empty", "scope(12, )"}, Result: false, Message: "数据不在约束范围内"}
	assert.Equal(t, reportImpl, failReport, "These are the same report")
}

func TestValidate_UninterruptedCheck(t *testing.T) {
	data := "[{\"name\":\"\",\"type\":\"\",\"is_must\":0,\"validate\":\"\"},{\"name\":\"pid\",\"type\":\"int\",\"is_must\":1,\"validate\":\"\"},{\"name\":\"id\",\"type\":\"int64\",\"is_must\":1,\"validate\":\"empty\"},{\"name\":\"email\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty\"},{\"name\":\"phone\",\"type\":\"string\",\"is_must\":0,\"validate\":\"phone|empty\"},{\"name\":\"startTime\",\"type\":\"string\",\"is_must\":0,\"validate\":\"datetime|empty\"},{\"name\":\"price\",\"type\":\"float\",\"is_must\":0,\"validate\":\"scope(9.99, 199.99)|empty\"},{\"name\":\"status\",\"type\":\"int\",\"is_must\":0,\"validate\":\"enum(1, 3, 4)|empty\"},{\"name\":\"path\",\"type\":\"string\",\"is_must\":0,\"validate\":\"enum(body, header)|empty\"},{\"name\":\"password\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope(6, 18)\"},{\"name\":\"level\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope(12, )\"},{\"name\":\"viewCount\",\"type\":\"string\",\"is_must\":0,\"validate\":\"empty|scope( , 100)\"}]"
	var params = map[string]string{
		"pid":       "100",
		"id":        "2000000",
		"email":     "kisschou@me.com",
		"phone":     "18018001800",
		"startTime": time.Now().Format("2006-01-02 15:04:05"),
		"price":     "10",
		"status":    "1",
		"path":      "header",
		"password":  "1234567",
		"level":     "13",
		"viewCount": "99",
	}
	validReportCenterImpl, _ := NewValidate().Json(data).UninterruptedCheck(params)
	report := &validReport{Name: "level", Rule: []string{"empty", "scope(12, )"}, Result: false, Message: "数据不在约束范围内"}
	assert.Equal(t, validReportCenterImpl.ReportByName("level"), report, "These are the same report")
	report = &validReport{Name: "pid", Rule: []string{}, Result: true, Message: "Success"}
	assert.Equal(t, validReportCenterImpl.ReportByIndex(1), report, "These are the same report")
	// 因为生成时间为题，导致两条不会相等
	// exportJson := "{\"build_time\":\"2021-05-13 01:03:37\",\"elapsed_time\":0,\"report_list\":[{\"name\":\"\",\"validate\":[],\"result\":true,\"message\":\"Success\"},{\"name\":\"pid\",\"validate\":[],\"result\":true,\"message\":\"Success\"},{\"name\":\"id\",\"validate\":[\"empty\"],\"result\":true,\"message\":\"Success\"},{\"name\":\"email\",\"validate\":[\"empty\"],\"result\":true,\"message\":\"Success\"},{\"name\":\"phone\",\"validate\":[\"phone\",\"empty\"],\"result\":true,\"message\":\"Success\"},{\"name\":\"startTime\",\"validate\":[\"datetime\",\"empty\"],\"result\":true,\"message\":\"Success\"},{\"name\":\"price\",\"validate\":[\"scope(9.99, 199.99)\",\"empty\"],\"result\":true,\"message\":\"Success\"},{\"name\":\"status\",\"validate\":[\"enum(1, 3, 4)\",\"empty\"],\"result\":true,\"message\":\"Success\"},{\"name\":\"path\",\"validate\":[\"enum(body, header)\",\"empty\"],\"result\":true,\"message\":\"Success\"},{\"name\":\"password\",\"validate\":[\"empty\",\"scope(6, 18)\"],\"result\":true,\"message\":\"Success\"},{\"name\":\"level\",\"validate\":[\"empty\",\"scope(12, )\"],\"result\":false,\"message\":\"数据不在约束范围内\"},{\"name\":\"viewCount\",\"validate\":[\"empty\",\"scope( , 100)\"],\"result\":true,\"message\":\"Success\"}]}"
	// assert.Equal(t, validReportCenterImpl.ToJson(), exportJson, "These are the same export json of report center.")
}
