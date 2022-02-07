package tdog

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	tc "github.com/kisschou/TypeConverter"
	"github.com/pelletier/go-toml"
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
		filePath    string // 配置文件路径
		defaultFile string // 默认配置文件
		fixedFile   string // 指定文件
		keyPrefix   string // key前缀
		searchKey   string // 查询的key值
		actionFile  string // 当前查询的文件名
		actionKey   string // 当前查询的key
	}

	// configResult result struct of configuration searched
	configResult struct {
		filePath   string      // 配置文件路径
		searchKey  string      // 查询的key值
		activeFile string      // 命中文件
		activeKey  string      // 命中key
		result     interface{} // 结果
		Message    string      // 消息
	}
)

// NewConfig init config struct
func NewConfig() *config {
	return &config{
		filePath:    os.Getenv("CONFIG_PATH"),
		defaultFile: "app",
		searchKey:   "",
		actionFile:  "",
		actionKey:   "",
	}
}

// connect Make go-toml load configuration given *config
func connect(c *config) *toml.Tree {
	conf, err := toml.LoadFile(c.filePath + c.actionFile + ".toml")
	if err != nil {
		Println(err.Error(), 13)
		os.Exit(0)
	}
	return conf
}

// SetPath assign configuration path
// It will change configuration path for search in linker
// given string path
// returns config-handler
func (c *config) SetPath(path string) *config {
	c.filePath = path
	return c
}

// SetFile assign and fixed configuration file name given string file name returns config-handler
func (c *config) SetFile(name string) *config {
	c.fixedFile = name
	return c
}

// SetPrefix assign and fixed key's prefix name.
// It will make your search key be prefix + key,
// So it will only be used with the GetMulti function.
// given string key's prefix
// returns config-handler
func (c *config) SetPrefix(prefix string) *config {
	c.keyPrefix = prefix
	return c
}

// find search from configuration returns *configResult and error struct when it err.
func (c *config) find() (*configResult, error) {
	var err error
	var result interface{}

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
		if (c.keyPrefix)[len(c.keyPrefix)-1:] != "." {
			c.keyPrefix += "."
		}
		c.actionKey = c.keyPrefix + c.actionKey
	}
	if len(c.filePath) < 1 {
		Println("Please set CONFIG_PATH in environment at first.", 13)
		os.Exit(0)
	} else {
		if (c.filePath)[len(c.filePath)-1:] != "/" {
			c.filePath += "/"
		}
		if runtime.GOOS == "windows" {
			c.filePath = strings.ReplaceAll(c.filePath, "/", "\\")
		}
	}

	for {
		configImpl := connect(c)
		if !configImpl.Has(c.actionKey) {
			// 固定文件的, 多段不应含有文件名
			if len(c.fixedFile) > 0 {
				err = errors.New(fmt.Sprintf("[%s.%s%s], 未找到配置, %s -> %s -> %s", c.fixedFile, c.keyPrefix, c.searchKey, c.filePath, c.actionFile, c.keyPrefix+c.actionKey))
				break
			}

			match := strings.Split(c.actionKey, ".")
			if len(match) > 1 {
				c.actionFile, c.actionKey = match[0], strings.Join(match[1:], ".")
				continue
			} else {
				err = errors.New(fmt.Sprintf("[%s], 未找到配置, %s -> %s -> %s", c.searchKey, c.filePath, c.actionFile, c.keyPrefix+c.actionKey))
				break
			}
		} else {
			result = configImpl.Get(c.actionKey)
			break
		}
	}

	return &configResult{
		filePath:   c.filePath,
		searchKey:  c.searchKey,
		activeFile: c.actionFile,
		activeKey:  c.actionKey,
		result:     result,
		Message:    "",
	}, err
}

// Get result by key on configuration file, extends *config
// given string the search you want
// returns *ConfigResult
func (c *config) Get(key string) *configResult {
	c.actionFile, c.actionKey, c.searchKey = "", "", key
	resultImpl, err := c.find()
	if err != nil {
		go NewLogger().Warn(err.Error())
		resultImpl.Message = err.Error()
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

// configResult -->

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
func (cr *configResult) RawData() interface{} {
	return cr.result
}

// ToString get the result of string type, if you sure about it, extends *configResult
// returns string
func (cr *configResult) ToString() string {
	return tc.New(cr.result).String
}

// ToInt get the result of int type, if you sure about it, extends *configResult
// returns int
func (cr *configResult) ToInt() int {
	return tc.New(cr.result).Int
}

// ToBool get the result of bool type, if you sure about it, extends *configResult
// returns bool
func (cr *configResult) ToBool() bool {
	return tc.New(cr.result).Bool
}

// ToIntSlice get the result of int slice type, if you sure about it, extends *configResult
// returns []int
func (cr *configResult) ToIntSlice() []int {
	s := make([]int, 0)
	is := tc.New(cr.result).Slice
	for _, v := range is {
		s = append(s, tc.New(v).Int)
	}
	return s
}

// ToStringMap get the result of string map type, if you sure about it, extends *configResult
// returns map[string]interface{}
func (cr *configResult) ToStringMap() map[string]interface{} {
	return tc.New(cr.result).Map
}

// ToStringMapString get the result of string map string type, if you sure about it, extends *configResult
// returns map[string]string
func (cr *configResult) ToStringMapString() map[string]string {
	m := make(map[string]string, 0)
	smi := tc.New(cr.result).Map
	for k, v := range smi {
		m[k] = tc.New(v).String
	}
	return m
}

// ToStringMapStringSlice get the result of string map string slice type, if you sure about it, extends *configResult
// returns map[string][]string
func (cr *configResult) ToStringMapStringSlice() map[string][]string {
	return (cr.result).(map[string][]string)
}

// ToStringSlice get the result of string slice type, if you sure about it, extends *configResult
// returns []string
func (cr *configResult) ToStringSlice() []string {
	s := make([]string, 0)
	is := tc.New(cr.result).Slice
	for _, v := range is {
		s = append(s, tc.New(v).String)
	}
	return s
}

// ToInt64 get the result of int64 type, if you sure about it, extends *configResult
// returns int64
func (cr *configResult) ToInt64() int64 {
	return tc.New(cr.result).Int64
}

// <--
