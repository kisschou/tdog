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

	// report 校验结果报告
	report struct {
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
	rules := make([]map[string]string, 0)
	err := json.Unmarshal([]byte(input), &rules)
	if err != nil {
		return v
	}
	UtilTdog := NewUtil()
	ruleList := make([]*Rule, 0)
	for _, ruleInfo := range rules {
		// 是否必须
		isMust, err := strconv.Atoi(ruleInfo["is_must"])
		if err != nil {
			continue
		}
		// 规则获取
		validateRule := strings.Split(ruleInfo["validate"], "|")
		for k, v := range validateRule {
			if !UtilTdog.InArray("[]string", v, ruleDict) {
				validateRule = UtilTdog.Remove("[]string", validateRule, k).([]string)
			}
		}
		ruleList = append(ruleList, &Rule{
			Name:      ruleInfo["name"],
			ParamType: ruleInfo["type"],
			IsMust:    isMust > 0,
			Rule:      validateRule,
		})
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
		if (err1 != nil && len(val.(string)) < min) || (err2 != nil && len(val.(string)) > max) {
			isSuccess = false
		}
		break
	case "int":
		if (err1 != nil && val.(int) < min) || (err2 != nil && val.(int) > max) {
			isSuccess = false
		}
		break
	case "int64":
		if (err1 != nil && val.(int64) < int64(min)) || (err2 != nil && val.(int64) > int64(max)) {
			isSuccess = false
		}
		break
	case "float":
		if (err1 != nil && val.(float32) < float32(min)) || (err2 != nil && val.(float32) > float32(max)) {
			isSuccess = false
		}
		break
	case "double":
		if (err1 != nil && val.(float64) < float64(min)) || (err2 != nil && val.(float64) > float64(max)) {
			isSuccess = false
		}
		break
	case "object":
		if (err1 != nil && len(val.(map[string]interface{})) < min) || (err2 != nil && len(val.(map[string]interface{})) > max) {
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

	if valType != "string" {
		err = errors.New("枚举只支持类型为string的数据")
		return
	}

	matchs := regexp.MustCompile(`^scope\((.*),(.*)\)$`).FindStringSubmatch(pattern)
	if len(matchs) != 2 {
		err = errors.New("规则中的枚举(enum)格式有问题.Example:enum(小黄,yellow,Mr.黄)")
		return
	}

	enums := strings.Split(matchs[1], ",")
	if NewUtil().InArray("[]string", val.(string), enums) {
		isSuccess = true
	}

	return
}

// checkIn 校验就是所有规则跑一遍
func checkIn(rule *Rule, needle map[string]string) (output *report, err error) {
	UtilTdog := NewUtil()
	defer Recover()

	if UtilTdog.Isset("map[string]string", rule.Name, needle) {
		var val interface{}
		val, err = verfityType(needle[rule.Name], rule.ParamType) // 值

		// 参数类型校验
		if err != nil {
			output = &report{Name: rule.Name, Rule: rule.Rule, Result: false, Message: "值类型与设定类型不符"}
			return
		}

		// 规则校验
		for _, ruleName := range rule.Rule {
			switch ruleName {
			case "empty": // 非空
				if UtilTdog.Empty("map[string]string", rule.Name, needle) {
					output = &report{Name: rule.Name, Rule: rule.Rule, Result: false, Message: "数据为空"}
				}
				break
			case "phone": // 手机号码
				if rule.ParamType != "string" || !UtilTdog.VerifyPhone(val.(string)) {
					output = &report{Name: rule.Name, Rule: rule.Rule, Result: false, Message: "号码格式错误"}
					return
				}
				break
			case "email": // 邮箱
				if rule.ParamType != "string" || !UtilTdog.VerifyEmail(val.(string)) {
					output = &report{Name: rule.Name, Rule: rule.Rule, Result: false, Message: "邮箱格式错误"}
					return
				}
				break
			case "date": // 日期
				if rule.ParamType != "string" || !UtilTdog.VerifyDate(val.(string)) {
					output = &report{Name: rule.Name, Rule: rule.Rule, Result: false, Message: "日期格式错误"}
					return
				}
				break
			case "datetime": // 日期时间
				if rule.ParamType != "string" || !UtilTdog.VerifyDateTime(val.(string)) {
					output = &report{Name: rule.Name, Rule: rule.Rule, Result: false, Message: "日期时间格式错误"}
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
						output = &report{Name: rule.Name, Rule: rule.Rule, Result: false, Message: "数据不在约束范围内"}
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
						output = &report{Name: rule.Name, Rule: rule.Rule, Result: false, Message: "数据不在枚举内"}
						return
					}
				}
				break
			}
		}
	} else if rule.IsMust {
		// 必填校验
		// 记录错误并跳出循环
		output = &report{Name: rule.Name, Rule: rule.Rule, Result: false, Message: "未包含"}
	}
	return
}

// Check 校验
// 一旦遇到校验失败的项, 立刻停止并返回报告.
func (v *validate) Check(needle map[string]string) (output *report, err error) {
	if len(v.rules) < 1 {
		err = errors.New("未指定校验规则")
		return
	}
	if len(needle) < 1 {
		err = errors.New("未发现需要校验的数据")
		return
	}
	for _, validateInfo := range v.rules {
		output, err = checkIn(validateInfo, needle)
		if err != nil {
			return
		}
		if !output.Result {
			return
		}
	}
	return
}

// UninterruptedCheck 无中断校验
// 遇到失败的项, 只记录，等所有数据都校验后,统一返回.
func (v *validate) UninterruptedCheck(needle map[string]string) (output []*report, err error) {
	if len(v.rules) < 1 {
		err = errors.New("未指定校验规则")
		return
	}
	if len(needle) < 1 {
		err = errors.New("未发现需要校验的数据")
		return
	}
	for _, validateInfo := range v.rules {
		eachOutput := new(report)
		eachOutput, err = checkIn(validateInfo, needle)
		if err != nil {
			return
		}
		output = append(output, eachOutput)
	}
	return
}
