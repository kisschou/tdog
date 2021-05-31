// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package config

type (
	tomlLexer struct {
		input []rune
		index int
		data  map[string]interface{}
	}

	tomlClass struct {
		startIndex int
		endIndex   int
		value      string
	}

	tomlTree struct {
	}
)

func newToml(input []rune) *tomlLexer {
	return &tomlLexer{
		input: input,
	}
}

func (tl *tomlLexer) leftBracket() {
}

func (tl *tomlLexer) rightBracket() {
}

func (tl *tomlLexer) next() rune {
}

func (tl *tomlLexer) leftSquareBrackets() {
}

func (tl *tomlLexer) rightSquareBrackets() {
}

func (tl *tomlLexer) leftBrace() {
}

func (tl *tomlLexer) rightBrace() {
}

func (tl *tomlLexer) dot() {
}

func (tl *tomlLexer) equal() {
}

func (tl *tomlLexer) commend() {
}

func (tl *tomlLexer) doubleQuotationMarks() {
}

func (tl *tomlLexer) singleQuotationMarks() {
}

func (tl *tomlLexer) comma() {
}

func (tl *tomlLexer) enter() {
}

func (tl *tomlLexer) lineBreak() {
}

func (tl *tomlLexer) query() {
	for {
		next := tl.next()
		switch next {
		case '(':
			break
		case ')':
			break
		case '[':
			break
		case ']':
			break
		case '{':
			break
		case '}':
			break
		case '.':
			break
		case '=':
			break
		case '#':
			break
		case '"':
			break
		case '\'':
			break
		case ',':
			break
		case '\r':
			break
		case '\n':
			break
		}
	}
}

func (tl *tomlLexer) run() {
	for status := tl.query; status != nil; {
	}
}
