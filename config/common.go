// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package config

import (
	"io/ioutil"
	"os"
	"strings"
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
