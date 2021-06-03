// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package config

type (
	// tomllexer .
	tomlLexer struct {
		input        []rune
		index        int
		currentTree  string
		currentLabel string
		box          *tomlTree
	}

	// tomlTree .
	tomlTree struct {
		childTree     map[string]*tomlTree
		landscapeList []map[string]string
	}
)

const (
	// monitorstart Mark the start of monitoring.
	monitorStop = iota
	// monitorstart Mark the stop of monitoring.
	monitorStart
)

var (
	// landscape .
	landscape []rune

	// squarebracketsMonitor Bracket listener status.
	squareBracketsMonitor int
	// squareBracketsSameCount The number of brackets monitored after the listener started.
	squareBracketsSameCount int
	// braceSameCount The number of brackets.
	braceSameCount int
	// equalMonitor Equal listener status.
	equalMonitor int
	// singleQuotationMarksSameCount The number of singleQuotationMarks monitored after the listener started.
	singleQuotationMarksSameCount int
	// singleQuotationMarksTripleMonitor SingleQuotationMarksTriple listener status.
	singleQuotationMarksTripleMonitor int
	// doubleQuotationMarksSameCount The number of doubleQuotationMarks monitored after the listener started.
	doubleQuotationMarksSameCount int
	// doubleQuotationMarksTripleMonitor DoubleQuotationMarksTripleMonitor listener status.
	doubleQuotationMarksTripleMonitor int
	// commendMonitor Commend listener status.
	commendMonitor int
)

// newtoml init a new toml lexer.
func newToml(input []rune) *tomlLexer {
	return &tomlLexer{
		input: input,
	}
}

// inBox check label is in tl.box.
func (tl *tomlLexer) inBox(label string) bool {
	if tl.box == nil {
		tl.box = new(tomlTree)
		return false
	}
	if _, ok := tl.box.childTree[label]; ok {
		return true
	}
	return false
}

// putCategory put the categray in tl.box.
func (tl *tomlLexer) putCategory() {
	// Convert to string and remove spaces.
	landscape = removeSpace(landscape)

	// set current tree's name.
	tl.currentTree = string(landscape)

	// put anything in box.
	if !tl.inBox(string(landscape)) {
		if tl.box.childTree == nil {
			tl.box.childTree = make(map[string]*tomlTree, 0)
		}
		tl.box.childTree[string(landscape)] = new(tomlTree)
	}

	// clear .
	landscape = make([]rune, 0)
}

// putInBox put k-v data in tl.box.
func (tl *tomlLexer) putInBox(landscape string) {
	if len(tl.currentLabel) < 1 {
		return
	}
	box := new(tomlTree)
	if tl.box != nil {
		box = tl.box
	}
	if len(tl.currentTree) < 1 {
		if box.landscapeList == nil {
			box.landscapeList = make([]map[string]string, 0)
		}
		box.landscapeList = append(box.landscapeList, map[string]string{tl.currentLabel: landscape})
	}
	if tl.inBox(tl.currentTree) {
		box.childTree[tl.currentTree].landscapeList = append(box.childTree[tl.currentTree].landscapeList, map[string]string{tl.currentLabel: landscape})
	}
	tl.box = box
}

// next get next character.
func (tl *tomlLexer) next() rune {
	if tl.index >= len(tl.input) {
		return -1
	}
	return tl.input[tl.index]
}

// leftSquareBrackets Dealing with opening left square brackets.
func (tl *tomlLexer) leftSquareBrackets() {
	if equalMonitor == monitorStart {
		squareBracketsSameCount++
		return
	}

	if squareBracketsMonitor == monitorStart {
		squareBracketsSameCount += 1
	} else {
		landscape = make([]rune, 0)
		squareBracketsMonitor = monitorStart
	}
}

// rightSquareBrackets Dealing with opening right square brackets.
func (tl *tomlLexer) rightSquareBrackets() {
	if equalMonitor == monitorStart {
		squareBracketsSameCount--
		return
	}
	if squareBracketsMonitor == monitorStart {
		if squareBracketsSameCount-1 >= 0 {
			squareBracketsSameCount--
		} else {
			// end SquareBrackets
			squareBracketsMonitor = monitorStop

			if landscape[len(landscape)-1] == ']' {
				landscape = landscape[:len(landscape)-1]
			}
			tl.putCategory()
		}
	}
}

// leftBrace Dealing with opening left brace.
func (tl *tomlLexer) leftBrace() {
	if equalMonitor == monitorStart {
		braceSameCount++
		return
	}
}

// rightBrace Dealing with opening right brace.
func (tl *tomlLexer) rightBrace() {
	if equalMonitor == monitorStart {
		braceSameCount--
		return
	}
}

// equal Dealing with opening equal.
func (tl *tomlLexer) equal() {
	if equalMonitor == monitorStart {
		equalMonitor = monitorStop
		if landscape[0] == '=' {
			landscape = landscape[1:]
		}
		tl.putInBox(string(removeSpace(landscape)))
	} else {
		equalMonitor = monitorStart
		if landscape[len(landscape)-1] == '=' {
			landscape = landscape[:len(landscape)-1]
		}
		tl.currentLabel = string(removeSpace(landscape))
		landscape = make([]rune, 0)
	}
}

// commend Dealing with commend.
func (tl *tomlLexer) commend() {
	if commendMonitor == monitorStart {
		commendMonitor = monitorStop
		landscape = make([]rune, 0)
	} else {
		commendMonitor = monitorStart
	}
}

// singleQuotationMarks Dealing with single quotation marks.
func (tl *tomlLexer) singleQuotationMarks() {
	// check is triple.
	if singleQuotationMarksSameCount == 1 {
		if tl.input[tl.index] == tl.input[tl.index-1] && tl.input[tl.index] == tl.input[tl.index+1] {
			singleQuotationMarksSameCount = -1
			if singleQuotationMarksTripleMonitor == monitorStart {
				singleQuotationMarksTripleMonitor = monitorStop
			} else {
				singleQuotationMarksTripleMonitor = monitorStart
			}
			return
		}
	}

	if equalMonitor == monitorStart {
		if singleQuotationMarksSameCount > 0 {
			singleQuotationMarksSameCount--
		} else {
			singleQuotationMarksSameCount++
		}
		return
	}
}

// doubleQuotationMarks Dealing with double quotation marks.
func (tl *tomlLexer) doubleQuotationMarks() {
	// check is triple.
	if doubleQuotationMarksSameCount == 1 {
		if tl.input[tl.index] == tl.input[tl.index-1] && tl.input[tl.index] == tl.input[tl.index+1] {
			singleQuotationMarksSameCount = -1
			if doubleQuotationMarksTripleMonitor == monitorStart {
				doubleQuotationMarksTripleMonitor = monitorStop
			} else {
				doubleQuotationMarksTripleMonitor = monitorStart
			}
			return
		}
	}

	if equalMonitor == monitorStart {
		if doubleQuotationMarksSameCount > 0 {
			doubleQuotationMarksSameCount--
		} else {
			doubleQuotationMarksSameCount++
		}
		return
	}
}

// lineBreak Dealing with double line break.
func (tl *tomlLexer) lineBreak() {
	landscape = landscape[:len(landscape)-1]
	if equalMonitor == monitorStart {
		if squareBracketsSameCount > 0 || braceSameCount > 0 ||
			doubleQuotationMarksSameCount > 0 || singleQuotationMarksSameCount > 0 ||
			doubleQuotationMarksTripleMonitor == monitorStart || singleQuotationMarksTripleMonitor == monitorStart {
			landscape = append(landscape, '\n')
		} else {
			tl.equal()
			landscape = make([]rune, 0)
		}
	}
}

// collation Start data sorting.
func (tl *tomlLexer) collation() *tomlLexer {
	for {
		next := tl.next()

		// end
		if next == -1 {
			break
		}

		landscape = append(landscape, next)

		switch next {
		case '(':
			// setCatch("LBracket", tl.index)
			break
		case ')':
			// setCatch("RBracket", tl.index)
			break
		case '[':
			tl.leftSquareBrackets()
			break
		case ']':
			tl.rightSquareBrackets()
			break
		case '{':
			tl.leftBrace()
			break
		case '}':
			tl.rightBrace()
			break
		case '.':
			// setCatch("Dot", tl.index)
			break
		case '=':
			if equalMonitor == monitorStop {
				tl.equal()
			}
			break
		case '#':
			tl.commend()
			break
		case '"':
			tl.doubleQuotationMarks()
			break
		case '\'':
			tl.singleQuotationMarks()
			break
		case ',':
			// setCatch("Comma", tl.index)
			break
		case '\r':
			// setCatch("Enter", tl.index)
			break
		case '\n':
			tl.lineBreak()
			break
		case ' ':
			// setCatch("Space", tl.index)
			break
		}

		tl.index++
	}
	return tl
}

func (tl *tomlLexer) output() {
}

func (tl *tomlLexer) run() {
	tl.collation().output()
}
