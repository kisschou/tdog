package tdog

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		filePath    string   // 配置文件路径
		configFiles []string // 所有的配置文件
		defaultFile string   // 默认配置文件
		fileSuffix  string   // 文件后缀
		searchKey   string   // 查询的key值
		actionFile  string   // 当前查询的文件名
		actionKey   string   // 当前查询的key
	}

	ConfigResult struct {
		filePath   string // 配置文件路径
		activeFile string // 命中文件
		activeKey  string // 命中key
		Message    string // 消息
	}
)

func NewConfig() *Config {
	config := &Config{
		filePath:    "",
		configFiles: nil,
		defaultFile: getDefaultFile(),
		fileSuffix:  getFileSuffix(),
		searchKey:   "",
		actionFile:  "",
		actionKey:   "",
	}
	config.SetPath(getFilePath())
	return config
}

func getFilePath() string {
	path := os.Getenv("CONFIG_PATH")
	if len(path) < 1 {
		path, _ = os.Getwd()
		path += "/config"
	}
	return path
}

func getDefaultFile() string {
	return "app"
}

func getFileSuffix() string {
	return "toml"
}

func (c *Config) getFiles() {
	if c.filePath == "" {
		c.configFiles = nil
	}
	c.configFiles, _ = NewUtil().GetFilesBySuffix(c.filePath, c.fileSuffix)
}

func (c *Config) SetPath(path string) *Config {
	c.filePath = path
	c.getFiles()
	return c
}

func connect(c *Config) {
	viper.SetConfigName(c.actionFile)
	viper.SetConfigType(c.fileSuffix)
	viper.AddConfigPath(c.filePath)
	err := viper.ReadInConfig()
	if err != nil {
		NewLogger().Error(err.Error())
	}
}

func (c *Config) find() (*ConfigResult, error) {
	var err error
	if len(c.actionFile) < 1 {
		c.actionFile = c.defaultFile
	}
	if len(c.actionKey) < 1 {
		c.actionKey = c.searchKey
	}
	for {
		connect(c)
		if !viper.IsSet(c.actionKey) {
			match := strings.Split(c.actionKey, ".")
			if len(match) > 1 {
				c.actionFile = match[0]
				if !NewUtil().InStringSlice(c.actionFile, c.configFiles) {
					err = errors.New(fmt.Sprintf("[%s], 无法定位到配置, %s -> %s", c.searchKey, c.filePath, c.actionFile))
					break
				}
				c.actionKey = strings.Join(match[1:], ".")
				if !viper.IsSet(c.actionKey) {
					err = errors.New(fmt.Sprintf("[%s], 未找到配置, %s -> %s -> %s", c.searchKey, c.filePath, c.actionFile, c.actionKey))
					break
				}
			}
		}
	}
	resultImpl := &ConfigResult{
		filePath:   c.filePath,
		activeFile: c.actionFile,
		activeKey:  c.actionKey,
		Message:    "",
	}
	return resultImpl, err
}

func (c *Config) Get(key string) *ConfigResult {
	c.searchKey = key
	resultImpl, err := c.find()
	if err != nil {
		resultImpl.Message = err.Error()
		return resultImpl
	}
	return resultImpl
}

func (cr *ConfigResult) IsExists() bool {
	isExists := false
	if len(cr.Message) < 1 {
		isExists = true
	}
	return isExists
}

func (c *ConfigResult) RawData() (data interface{}) {
	connect(&Config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.Get(c.activeKey)
	return
}

func (c *ConfigResult) String() (data string) {
	connect(&Config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetString(c.activeKey)
	return
}

func (c *ConfigResult) Int() (data int) {
	connect(&Config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetInt(c.activeKey)
	return
}

func (c *ConfigResult) Bool() (data bool) {
	connect(&Config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetBool(c.activeKey)
	return
}

func (c *ConfigResult) IntSlice() (data []int) {
	connect(&Config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetIntSlice(c.activeKey)
	return
}

func (c *ConfigResult) StringMap() (data map[string]interface{}) {
	connect(&Config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringMap(c.activeKey)
	return
}

func (c *ConfigResult) StringMapString() (data map[string]string) {
	connect(&Config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringMapString(c.activeKey)
	return
}

func (c *ConfigResult) StringMapStringSlice() (data map[string][]string) {
	connect(&Config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringMapStringSlice(c.activeKey)
	return
}

func (c *ConfigResult) StringSlice() (data []string) {
	connect(&Config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringSlice(c.activeKey)
	return
}

func (c *ConfigResult) Int64() (data int64) {
	connect(&Config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetInt64(c.activeKey)
	return
}
