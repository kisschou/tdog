// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

type (
	box struct {
		childBox      map[string]*box
		landscapeList map[string][]landscapeInfo
	}

	landscapeInfo struct {
		snippetType int
		snippet     interface{}
	}
)

const (
	typeInterface = iota
	typeString
	typeInt
	typeInt64
	typeBool

	typeSliceInterface
	typeSliceString
	typeSliceInt
	typeSliceInt64
	typeSliceMapStringString
	typeSliceMapStringSliceString
	typeSliceMapStringSliceInt
	typeSliceMapStringSliceInt64
	typeSliceMapStringSliceInterface

	typeMapStringInterface
	typeMapStringString
	typeMapStringInt
	typeMapStringInt64
	typeMapStringSliceString
	typeMapStringSliceInt
	typeMapStringSliceInt64
)

// getFilesBySuffix Gets the filename of all the specified suffixes from the specified path.
// given string path means file path of scan
// given string suffix means catch for same suffix
// returns []string files list of file name, file name has no suffix
// returns error err throw it if has errors
func getFilesBySuffix(path string, suffix string) (files []string, err error) {
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			fileSuffix := path.Ext(fi.Name())
			fileName := strings.TrimRight(path, '/') + "/" + fi.Name()
			if "."+suffix == fileSuffix {
				files = append(files, fileName)
			}
		}
	}
	return
}

// fileExists check file is exists. given string path returns true when exists
func fileExists(input string) bool {
	_, err := os.Stat(input)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// isDir check path is dir given string path returns true when it's
func isDir(input string) bool {
	s, err := os.Stat(input)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// isFile check path is file given string path returns true when it's
func isFile(input string) bool {
	if !isDir(input) {
		return fileExists(input)
	}
	return false
}

// getContent get content of file.
func getContent(input string) (data []byte, err error) {
	data, err = ioutil.ReadFile(input)
	return
}

// removeSpace like its name. But only delete the leading and trailing spaces.
func removeSpace(input []rune) []rune {
	res := input
	for k, v := range input {
		if v == ' ' {
			res = input[k+1:]
		} else {
			break
		}
	}
	input = res
	for i := len(input) - 1; i > 0; i-- {
		if input[i] == ' ' {
			res = input[:i]
		} else {
			break
		}
	}
	return res
}

func snippetTypeLexel(actualType int, actualSnippet interface{}) []landscapeInfo {
}

func snippetAnalysis(tree *tomlTree) (box *box) {
	if tree.childTree != nil {
		for categoryName, childTree := range tree.childTree {
			box.childBox[categoryName] = snippetAnalysis(childTree)
		}
	}
	if tree.landscapeList != nil {
		for label, snippet := range tree.landscapeList {
			box.landscapeList[label] = snippetLexel(bytes.Runes([]byte(snippet)))
		}
	}
	return
}

func (tl *tomlLexer) query() *box {
	return snippetAnalysis(tl.box)
}
