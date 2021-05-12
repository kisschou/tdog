// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/**
 * The module for validate handler.
 *
 * Author: Kisschou
 * @Build: 2021-05-11
 */
type (
	// validate 自动验证模块
	validate struct {
		rules []*Rule
	}

	// Rule 校验规则
	Rule struct {
		Name      string   `json:"name"`     // 名称
		ParamType string   `json:"type"`     // 类型
		IsMust    bool     `json:"is_must"`  // 是否必须
		Rule      []string `json:"validate"` // 规则
	}

	// reportCenter 校验报告中心
	validReportCenter struct {
		reports     []*validReport `json:"report_list"`  // 校验报告列表
		createTime  string         `json:"build_time"`   // 报告生成时间
		elapsedTime int64          `json:"elapsed_time"` // 执行耗时
	}

	// report 校验结果报告
	validReport struct {
		Name    string   `json:"name"`     // 名称
		Rule    []string `json:"validate"` // 校验方式
		Result  bool     `json:"result"`   // 校验结果
		Message string   `json:"message"`  // 结果信息
	}
)

var (
	// ruleDict 校验规则字典
	ruleDict = []string{
		"phone",          // 手机号码
		"email",          // 邮箱
		"empty",          // 非空
		"scope",          // 长度及取值范围 // 类型为字符串时判断为字符串长度范围 // example: scope(0,): 0<x; scope(0,10): 0<x<10; scope(,10): x<10
		"enum",           // 枚举 // example: enum(0,1,3); enum(小黄,yellow,Mr.黄);
		"date",           // 日期 // 标准格式: yyyy-mm-dd
		"datetime",       // 日期时间 // 标准格式: yyyy-mm-dd hh:mm:ss
		"sensitive-word", // 敏感词 // 外部关联，暂不支持
	}
)

// NewValidate init a new validate model
func NewValidate() *validate {
	return &validate{}
}

// Rule Set a rule that you want to use in the next validation.
// given rule extend Rule struct, returns validate struct.
func (v *validate) Rule(input []*Rule) *validate {
	v.rules = input
	return v
}

// Json As with Rule, the validation rules are set.
// The difference is that this is of JSON type and needs to be converted to []*Rule.
// example: [{"name":"api_id","type":"int64","is_must":1,"validate":"scope(0,)|"},...]
func (v *validate) Json(input string) *validate {
	rules := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(input), &rules)
	if err != nil {
		return v
	}
	UtilTdog := NewUtil()
	ruleList := make([]*Rule, 0)
	for _, ruleInfo := range rules {
		// 是否必须
		isMust := int(ruleInfo["is_must"].(float64))
		// 规则获取
		validateRule := strings.Split(ruleInfo["validate"].(string), "|")
		for k, v := range validateRule {
			validateRule[k] = strings.TrimSpace(v)
			if len(strings.TrimSpace(v)) < 1 {
				validateRule = UtilTdog.Remove("[]string", validateRule, k).([]string)
			}
		}
		eachRule := &Rule{
			Name:      ruleInfo["name"].(string),
			ParamType: ruleInfo["type"].(string),
			IsMust:    isMust > 0,
			Rule:      validateRule,
		}
		ruleList = append(ruleList, eachRule)
	}
	v.rules = ruleList
	return v
}

// verifyType 字符串类型校验
func verfityType(input, actionType string) (res interface{}, err error) {
	switch actionType {
	default:
	case "string":
		res = input
		break
	case "int":
		res, err = strconv.Atoi(input)
		break
	case "int64":
		res, err = strconv.ParseInt(input, 10, 64)
		break
	case "float":
		res, err = strconv.ParseFloat(input, 32)
		break
	case "double":
		res, err = strconv.ParseFloat(input, 64)
		break
	case "object":
		resMap := make(map[string]interface{}, 0)
		err = json.Unmarshal([]byte(input), &resMap)
		res = resMap
		break
	}
	return
}

// verifyScore 校验是否符合范围约束
func verifyScore(pattern, valType string, val interface{}) (isSuccess bool, err error) {
	defer Recover()

	matchs := regexp.MustCompile(`^scope\((.*),(.*)\)$`).FindStringSubmatch(pattern)
	if len(matchs) != 3 {
		err = errors.New("规则中的取值范围(scope)格式有问题.Example:scope(0,20)")
		return
	}
	min, err1 := strconv.Atoi(strings.TrimSpace(matchs[1]))
	max, err2 := strconv.Atoi(strings.TrimSpace(matchs[2]))
	isSuccess = true
	switch valType {
	default:
		break
	case "string":
		if (err1 == nil && len(val.(string)) < min) || (err2 == nil && len(val.(string)) > max) {
			isSuccess = false
		}
		break
	case "int":
		if (err1 == nil && val.(int) < min) || (err2 == nil && val.(int) > max) {
			isSuccess = false
		}
		break
	case "int64":
		min, err1 := strconv.ParseInt(strings.TrimSpace(matchs[1]), 10, 64)
		max, err2 := strconv.ParseInt(strings.TrimSpace(matchs[2]), 10, 64)
		if (err1 == nil && val.(int64) < min) || (err2 == nil && val.(int64) > max) {
			isSuccess = false
		}
		break
	case "float", "double":
		min, err1 := strconv.ParseFloat(strings.TrimSpace(matchs[1]), 64)
		max, err2 := strconv.ParseFloat(strings.TrimSpace(matchs[2]), 64)
		if (err1 == nil && val.(float64) < min) || (err2 == nil && val.(float64) > max) {
			isSuccess = false
		}
		break
	case "object":
		if (err1 == nil && len(val.(map[string]interface{})) < min) || (err2 == nil && len(val.(map[string]interface{})) > max) {
			isSuccess = false
		}
		break
	}
	return
}

// verifEnum 校验是否符合枚举约束
func verifyEnum(pattern, valType string, val interface{}) (isSuccess bool, err error) {
	defer Recover()
	isSuccess = false

	if !NewUtil().InArray("[]string", valType, []string{"string", "int"}) {
		err = errors.New("枚举只支持类型为string和int的数据")
		return
	}

	matchs := regexp.MustCompile(`^enum\((.*)\)$`).FindStringSubmatch(pattern)
	if len(matchs) != 2 {
		err = errors.New("规则中的枚举(enum)格式有问题.Example:enum(小黄,yellow,Mr.黄)")
		return
	}

	enums := strings.Split(matchs[1], ",")
	for k, v := range enums {
		enums[k] = strings.TrimSpace(v)
	}
	switch valType {
	case "string":
		if NewUtil().InArray("[]string", val.(string), enums) {
			isSuccess = true
		}
		break
	case "int":
		if NewUtil().InArray("[]string", strconv.Itoa(val.(int)), enums) {
			isSuccess = true
		}
		break
	}

	return
}

// checkIn 校验就是所有规则跑一遍
func checkIn(rule *Rule, needle map[string]string) (err error) {
	UtilTdog := NewUtil()
	defer Recover()

	if UtilTdog.Isset("map[string]string", rule.Name, needle) {
		var val interface{}
		val, err = verfityType(needle[rule.Name], rule.ParamType) // 值

		// 参数类型校验
		if err != nil {
			err = errors.New("值类型与设定类型不符")
			return
		}

		// 规则校验
		for _, ruleName := range rule.Rule {
			switch ruleName {
			case "empty": // 非空
				if UtilTdog.Empty("map[string]string", rule.Name, needle) {
					err = errors.New("数据为空")
					return
				}
				break
			case "phone": // 手机号码
				if rule.ParamType != "string" || !UtilTdog.VerifyPhone(val.(string)) {
					err = errors.New("号码格式错误")
					return
				}
				break
			case "email": // 邮箱
				if rule.ParamType != "string" || !UtilTdog.VerifyEmail(val.(string)) {
					err = errors.New("邮箱格式错误")
					return
				}
				break
			case "date": // 日期
				if rule.ParamType != "string" || !UtilTdog.VerifyDate(val.(string)) {
					err = errors.New("日期格式错误")
					return
				}
				break
			case "datetime": // 日期时间
				if rule.ParamType != "string" || !UtilTdog.VerifyDateTime(val.(string)) {
					err = errors.New("日期时间格式错误")
					return
				}
				break
			case "sensitive-word": // 敏感词
				break
			default:
				// 范围
				if strings.Contains(ruleName, "scope") {
					var isSuccess bool
					isSuccess, err = verifyScore(ruleName, rule.ParamType, val)
					if err != nil {
						return
					}
					if !isSuccess {
						err = errors.New("数据不在约束范围内")
						return
					}
				}
				// 枚举
				if strings.Contains(ruleName, "enum") {
					var isSuccess bool
					isSuccess, err = verifyEnum(ruleName, rule.ParamType, val)
					if err != nil {
						return
					}
					if !isSuccess {
						err = errors.New("数据不在枚举内")
						return
					}
				}
				break
			}
		}
	} else if rule.IsMust {
		// 必填校验
		// 记录错误并跳出循环
		err = errors.New("未包含")
		return
	}

	err = nil
	return
}

// Check 校验
// 一旦遇到校验失败的项, 立刻停止并返回报告.
func (v *validate) Check(needle map[string]string) (output *validReport, err error) {
	if len(v.rules) < 1 {
		err = errors.New("未指定校验规则")
		return
	}
	if len(needle) < 1 {
		err = errors.New("未发现需要校验的数据")
		return
	}
	for _, validateInfo := range v.rules {
		output = &validReport{Name: validateInfo.Name, Rule: validateInfo.Rule, Result: true, Message: "Success"}
		err = checkIn(validateInfo, needle)
		if err != nil {
			output.Result = false
			output.Message = err.Error()
			err = nil
			return
		}
	}
	// all success.
	output = &validReport{Name: "", Rule: []string{}, Result: true, Message: "Success"}
	return
}

// UninterruptedCheck 无中断校验
// 遇到失败的项, 只记录，等所有数据都校验后,统一返回.
// 返回报告中心结构体
func (v *validate) UninterruptedCheck(needle map[string]string) (output *validReportCenter, err error) {
	if len(v.rules) < 1 {
		err = errors.New("未指定校验规则")
		return
	}
	if len(needle) < 1 {
		err = errors.New("未发现需要校验的数据")
		return
	}
	s := time.Now().UnixNano()
	output = new(validReportCenter)
	reports := make([]*validReport, 0)
	for _, validateInfo := range v.rules {
		report := &validReport{Name: validateInfo.Name, Rule: validateInfo.Rule, Result: true, Message: "Success"}
		err = checkIn(validateInfo, needle)
		if err != nil {
			report.Result = false
			report.Message = err.Error()
		}
		reports = append(reports, report)
	}
	err = nil
	output.reports = reports
	output.createTime = time.Now().Format("2006-01-02 15:04:05")
	output.elapsedTime = (time.Now().UnixNano() - s) / 1000000
	return
}

// ReportList get all report from report center.
func (rc *validReportCenter) ReportList() []*validReport {
	return rc.reports
}

// ReportByIndex get the report by index. so given int index, returns *report.
func (rc *validReportCenter) ReportByIndex(index int) *validReport {
	return rc.reports[index]
}

// ReportByName get the report by name. so must given string name, and will returns *report
func (rc *validReportCenter) ReportByName(name string) *validReport {
	for _, reportImpl := range rc.reports {
		if reportImpl.Name == name {
			return reportImpl
		}
	}
	return &validReport{Name: "", Rule: []string{}, Result: false, Message: "未找到名为`" + name + "`的报表."}
}

// BuildTime get build time from report center.
func (rc *validReportCenter) BuildTime() string {
	return rc.createTime
}

// ElapsedTime get elapsed time from report center.
func (rc *validReportCenter) ElapsedTime() int64 {
	return rc.elapsedTime
}

// ToJson convert to json and return.
func (rc *validReportCenter) ToJson() string {
	var inputData = map[string]interface{}{
		"build_time":   rc.BuildTime(),
		"elapsed_time": rc.elapsedTime,
		"report_list":  rc.reports,
	}
	data, err := json.Marshal(inputData)
	if err != nil {
		go NewLogger().Error(err.Error())
		return ""
	}
	return string(data)
}
