// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import (
	"encoding/json"
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
	}
)

var (
	// ruleDict 校验规则字典
	ruleDict = []string{
		"phone",          // 手机号码
		"email",          // 邮箱
		"empty",          // 非空
		"scope",          // 长度及取值范围 // 类型为字符串时判断为字符串长度范围 // example: scope(0,): 0<x; scope(0,10): 0<x<10; scope(,10): x<10
		"enum",           // 枚举 // example: enum(0,1,3); enum("小黄", "yellow", "Mr.黄");
		"date",           // 日期 // 标准格式: yyyy-mm-dd
		"datetime",       // 日期时间 // 标准格式: yyyy-mm-dd hh:mm:ss
		"sensitive-word", // 敏感词 // 外部关联，暂不支持
	}
)

// NewValidate init a new validate model
func NewValidate() *validate {
	return &validate{}
}

// defaultRule Get the default rule.
// The default rule is: must exist and cannot be empty.
func (v *validate) defaultRule(rule *Rule) *Rule {
	rule.IsMust = true
	rule.Rule = []string{"empty"}
	return rule
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
	ruleList := make([]*Rule, 0)
	for _, ruleInfo := range rules {
		// 是否必须
		isMust, err := strconv.Atoi(ruleInfo["is_must"])
		if err != nil {
			continue
		}
		// 规则获取
		UtilTdog := NewUtil()
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

// Check 校验
// 一旦遇到校验失败的项, 立刻停止并返回报告.
func (v *validate) Check(input map[string]string) {
}

// UninterruptedCheck 无中断校验
// 遇到失败的项, 只记录，等所有数据都校验后,统一返回.
func (v *validate) UninterruptedCheck(input map[string]string) {
}
