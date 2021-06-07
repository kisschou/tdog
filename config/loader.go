// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package config

import (
	"strings"
)

type ()

func Load(input string) {
	if strings.Contains(input, "=") {
		loadBytes([]byte(input))
	}
	s := strings.Split(input, ".")
	if len(s) > 1 && s[len(s)-1] == "toml" {
		loadFile(input)
	}
	loadDir(input)
}

func loadBytes(input []byte) {
	lexToml(input)
}

func loadDir(input string) {
	if isDir(input) {
		files, _ := getFilesBySuffix(input, "toml")
		for _, file := range files {
			loadFile(file)
		}
	}
}

func loadFile(input string) {
	if isFile(input) {
		data, err := getContent(input)
		if err != nil {
			return
		}
		loadBytes(data)
	}
}
