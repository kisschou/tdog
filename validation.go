// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import ()

/**
 * The module for validate handler.
 *
 * Author: Kisschou
 * @Build: 2021-05-11
 */
type (
	// validate 自动验证模块
	validate struct {
		rule *Rule
	}

	// Rule 校验规则
	Rule struct {
	}

	// report 校验结果报告
	report struct {
	}
)

// NewValidate init a new validate model
func NewValidate() *validate {
	return &validate{}
}

// Rule Set a rule that you want to use in the next validation.
// given rule extend Rule struct, returns validate struct.
func (v *validate) Rule(rule *Rule) *validate {
	v.rule = rule
	return v
}

func (v *validate) Check() {
}

func (v *validate) UninterruptedCheck() {
}
