// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

/**
 * The module for configuration handler.
 *
 * @Author: Kisschou
 * @Build: 2021-04-21
 */
type (
	// config configuration attributes
	config struct {
		filePath    string   // 配置文件路径
		configFiles []string // 所有的配置文件
		defaultFile string   // 默认配置文件
		fileSuffix  string   // 文件后缀
		fixedFile   string   // 指定文件
		keyPrefix   string   // key前缀
		searchKey   string   // 查询的key值
		actionFile  string   // 当前查询的文件名
		actionKey   string   // 当前查询的key
	}

	// configResult result struct of configuration searched
	configResult struct {
		filePath   string // 配置文件路径
		searchKey  string // 查询的key值
		activeFile string // 命中文件
		activeKey  string // 命中key
		Message    string // 消息
	}
)

// NewConfig init config struct
func NewConfig() *config {
	configTdog := &config{
		filePath:    "",
		configFiles: nil,
		defaultFile: getDefaultFile(),
		fileSuffix:  getFileSuffix(),
		searchKey:   "",
		actionFile:  "",
		actionKey:   "",
	}
	configTdog.SetPath(getFilePath())
	return configTdog
}

// getFilePath returns defined configuration path when "CONFIG_PATH" isset in environment. default return /path/to/config
func getFilePath() string {
	path := os.Getenv("CONFIG_PATH")
	if len(path) < 1 {
		path, _ = os.Getwd()
		path += "/config"
	}
	return path
}

// getFilePath returns default configuration file name
func getDefaultFile() string {
	return "app"
}

// getFilePath returns default configuration file suffix
func getFileSuffix() string {
	return "toml"
}

// getFilePath get all configuration file name from path
func (c *config) getFiles() {
	if c.filePath == "" {
		c.configFiles = nil
	}
	c.configFiles, _ = NewUtil().GetFilesBySuffix(c.filePath, c.fileSuffix)
}

// SetPath assign configuration path
// It will change configuration path for search in linker
// given string path
// returns config-handler
func (c *config) SetPath(path string) *config {
	c.filePath = path
	c.getFiles()
	return c
}

// SetFile assign and fixed configuration file name given string file name returns config-handler
func (c *config) SetFile(name string) *config {
	c.fixedFile = name
	return c
}

// SetPrefix assign and fixed key's prefix name
// It will make your search key be prefix + key
// given string key's prefix
// returns config-handler
func (c *config) SetPrefix(prefix string) *config {
	c.keyPrefix = prefix
	return c
}

// connect Make viper load configuration given *config
func connect(c *config) {
	viper.SetConfigName(c.actionFile)
	viper.SetConfigType(c.fileSuffix)
	viper.AddConfigPath(c.filePath)
	err := viper.ReadInConfig()
	if err != nil {
		go NewLogger().Error(err.Error())
	}
}

// find search from configuration returns *configResult and error struct when it err
func (c *config) find() (*configResult, error) {
	var err error
	if len(c.actionFile) < 1 {
		c.actionFile = c.defaultFile
	}
	if len(c.fixedFile) > 0 {
		c.actionFile = c.fixedFile
	}
	if len(c.actionKey) < 1 {
		c.actionKey = c.searchKey
	}
	if len(c.keyPrefix) > 0 {
		c.actionKey = c.keyPrefix + c.actionKey
	}
	for {
		connect(c)
		if !viper.IsSet(c.actionKey) {
			if len(c.fixedFile) > 0 {
				err = errors.New(fmt.Sprintf("[%s.%s%s], 未找到配置, %s -> %s -> %s", c.fixedFile, c.keyPrefix, c.searchKey, c.filePath, c.actionFile, c.keyPrefix+c.actionKey))
				break
			}
			match := strings.Split(c.actionKey, ".")
			if len(match) > 1 {
				c.actionFile = match[0]
				if !NewUtil().InStringSlice(c.actionFile, c.configFiles) {
					err = errors.New(fmt.Sprintf("[%s], 无法定位到配置, %s -> %s", c.searchKey, c.filePath, c.actionFile))
					break
				}
				c.actionKey = strings.Join(match[1:], ".")
			} else {
				err = errors.New(fmt.Sprintf("[%s], 未找到配置, %s -> %s -> %s", c.searchKey, c.filePath, c.actionFile, c.keyPrefix+c.actionKey))
				break
			}
		} else {
			break
		}
	}
	resultImpl := &configResult{
		filePath:   c.filePath,
		searchKey:  c.searchKey,
		activeFile: c.actionFile,
		activeKey:  c.actionKey,
		Message:    "",
	}
	return resultImpl, err
}

// Get get result by key on configuration file, extends *config
// given string the search you want
// returns *ConfigResult
func (c *config) Get(key string) *configResult {
	c.actionFile, c.actionKey, c.searchKey = "", "", key
	c.searchKey = key
	resultImpl, err := c.find()
	if err != nil {
		go NewLogger().Warn(err.Error())
		resultImpl.Message = err.Error()
		return resultImpl
	}
	return resultImpl
}

// GetMulti query multiple results in batches
// given ...string all the search keys you want
// returns map[string]*ConfigResult means search key as key and *ConfigResult as value
func (c *config) GetMulti(keys ...string) map[string]*configResult {
	if len(keys) < 1 {
		go NewLogger().Warn("Config: 批量查询参数缺失.")
		return nil
	}
	multiConfigResult := make(map[string]*configResult, 0)
	for _, key := range keys {
		configResultImpl := c.Get(key)
		multiConfigResult[configResultImpl.searchKey] = configResultImpl
	}
	return multiConfigResult
}

// GetSearchKey get search key from search result, extends *configResult returns string search key
func (cr *configResult) GetSearchKey() string {
	return cr.searchKey
}

// IsExists check is got value, extends *configResult return bool true when exists
func (cr *configResult) IsExists() bool {
	isExists := false
	if len(cr.Message) < 1 {
		isExists = true
	}
	return isExists
}

// RawData get the result of interface type, extends *configResult
// returns interface
// can get true result by use x.(type)
func (c *configResult) RawData() (data interface{}) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.Get(c.activeKey)
	return
}

// ToString get the result of string type, if you sure about it, extends *configResult
// returns string
func (c *configResult) ToString() (data string) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetString(c.activeKey)
	return
}

// ToInt get the result of int type, if you sure about it, extends *configResult
// returns int
func (c *configResult) ToInt() (data int) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetInt(c.activeKey)
	return
}

// ToBool get the result of bool type, if you sure about it, extends *configResult
// returns bool
func (c *configResult) ToBool() (data bool) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetBool(c.activeKey)
	return
}

// ToIntSlice get the result of int slice type, if you sure about it, extends *configResult
// returns []int
func (c *configResult) ToIntSlice() (data []int) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetIntSlice(c.activeKey)
	return
}

// ToStringMap get the result of string map type, if you sure about it, extends *configResult
// returns map[string]interface{}
func (c *configResult) ToStringMap() (data map[string]interface{}) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringMap(c.activeKey)
	return
}

// ToStringMapString get the result of string map string type, if you sure about it, extends *configResult
// returns map[string]string
func (c *configResult) ToStringMapString() (data map[string]string) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringMapString(c.activeKey)
	return
}

// ToStringMapStringSlice get the result of string map string slice type, if you sure about it, extends *configResult
// returns map[string][]string
func (c *configResult) ToStringMapStringSlice() (data map[string][]string) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringMapStringSlice(c.activeKey)
	return
}

// ToStringSlice get the result of string slice type, if you sure about it, extends *configResult
// returns []string
func (c *configResult) ToStringSlice() (data []string) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringSlice(c.activeKey)
	return
}

// ToInt64 get the result of int64 type, if you sure about it, extends *configResult
// returns int64
func (c *configResult) ToInt64() (data int64) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetInt64(c.activeKey)
	return
}
