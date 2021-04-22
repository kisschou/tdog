package tdog

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

/**
 * 配置文件获取模块
 *
 * @Author: Kisschou
 * @Build: 2021-04-21
 */
type (
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

	configResult struct {
		filePath   string // 配置文件路径
		searchKey  string // 查询的key值
		activeFile string // 命中文件
		activeKey  string // 命中key
		Message    string // 消息
	}
)

/**
 * 初始化配置模块
 *
 * @return *Config
 */
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

/**
 * 获取默认文件检索位置
 *
 * @return string
 */
func getFilePath() string {
	path := os.Getenv("CONFIG_PATH")
	if len(path) < 1 {
		path, _ = os.Getwd()
		path += "/config"
	}
	return path
}

/**
 * 获取默认配置文件
 *
 * @return string
 */
func getDefaultFile() string {
	return "app"
}

/**
 * 获取默认文件后缀
 *
 * @return string
 */
func getFileSuffix() string {
	return "toml"
}

/**
 * 获取路径下的所有指定格式的配置文件名
 *
 * @return nil
 */
func (c *config) getFiles() {
	if c.filePath == "" {
		c.configFiles = nil
	}
	c.configFiles, _ = NewUtil().GetFilesBySuffix(c.filePath, c.fileSuffix)
}

/**
 * 指定配置检索路径
 *
 * @param string path 路径
 *
 * @return *Config
 */
func (c *config) SetPath(path string) *config {
	c.filePath = path
	c.getFiles()
	return c
}

/**
 * 指定检索的文件
 *
 * @param string name 文件名
 *
 * @return *Config
 */
func (c *config) SetFile(name string) *config {
	c.fixedFile = name
	return c
}

/**
 * 指定检索的key前缀
 *
 * @param string prefix 前缀
 *
 * @return *Config
 */
func (c *config) SetPrefix(prefix string) *config {
	c.keyPrefix = prefix
	return c
}

/**
 * 连接到指定配置文件
 *
 * @param *Config c 指定配置文件相关结构体
 *
 * @return nil
 */
func connect(c *config) {
	viper.SetConfigName(c.actionFile)
	viper.SetConfigType(c.fileSuffix)
	viper.AddConfigPath(c.filePath)
	err := viper.ReadInConfig()
	if err != nil {
		NewLogger().Error(err.Error())
	}
}

/**
 * 按照配置规则检索配置文件
 *
 * @return *ConfigResult 检索结果结构体
 * @return error 错误信息
 */
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

/**
 * 获取配置查询返回结果结构体
 *
 * @param string key 需要查询的key值
 *
 * @return *ConfigResult 结果结构体
 */
func (c *config) Get(key string) *configResult {
	c.actionFile, c.actionKey, c.searchKey = "", "", key
	c.searchKey = key
	resultImpl, err := c.find()
	if err != nil {
		NewLogger().Warn(err.Error())
		resultImpl.Message = err.Error()
		return resultImpl
	}
	return resultImpl
}

/**
 * 获取配置查询返回结果结构体
 *
 * @param ...string keys (可变参数)需要查询的key
 *
 * @return []*ConfigResult 结果结构体切片
 */
func (c *config) GetMulti(keys ...string) map[string]*configResult {
	if len(keys) < 1 {
		NewLogger().Warn("Config: 批量查询参数缺失.")
		return nil
	}
	multiConfigResult := make(map[string]*configResult, 0)
	for _, key := range keys {
		configResultImpl := c.Get(key)
		multiConfigResult[configResultImpl.searchKey] = configResultImpl
	}
	return multiConfigResult
}

/**
 * 获取结果对应的原始查询键
 *
 * @return string
 */
func (cr *configResult) GetSearchKey() string {
	return cr.searchKey
}

/**
 * 判断是否查询到结果
 *
 * @return bool
 */
func (cr *configResult) IsExists() bool {
	isExists := false
	if len(cr.Message) < 1 {
		isExists = true
	}
	return isExists
}

/**
 * 使用interface{}类型获取数据结果
 *
 * @return interface{}
 */
func (c *configResult) RawData() (data interface{}) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.Get(c.activeKey)
	return
}

/**
 * 使用String类型获取数据结果
 *
 * @return string
 */
func (c *configResult) ToString() (data string) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetString(c.activeKey)
	return
}

/**
 * 使用int类型获取数据结果
 *
 * @return int
 */
func (c *configResult) ToInt() (data int) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetInt(c.activeKey)
	return
}

/**
 * 使用Bool类型获取数据结果
 *
 * @return bool
 */
func (c *configResult) ToBool() (data bool) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetBool(c.activeKey)
	return
}

/**
 * 使用IntSlice类型获取数据结果
 *
 * @return []int
 */
func (c *configResult) ToIntSlice() (data []int) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetIntSlice(c.activeKey)
	return
}

/**
 * 使用StringMap类型获取数据结果
 *
 * @return map[string]interface{}
 */
func (c *configResult) ToStringMap() (data map[string]interface{}) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringMap(c.activeKey)
	return
}

/**
 * 使用StringMapString类型获取数据结果
 *
 * @return map[string]string
 */
func (c *configResult) ToStringMapString() (data map[string]string) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringMapString(c.activeKey)
	return
}

/**
 * 使用StringMapStringSlice类型获取数据结果
 *
 * @return map[string][]string
 */
func (c *configResult) ToStringMapStringSlice() (data map[string][]string) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringMapStringSlice(c.activeKey)
	return
}

/**
 * 使用StringSlice类型获取数据结果
 *
 * @return []string
 */
func (c *configResult) ToStringSlice() (data []string) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetStringSlice(c.activeKey)
	return
}

/**
 * 使用int64类型获取数据结果
 *
 * @return int64
 */
func (c *configResult) ToInt64() (data int64) {
	connect(&config{filePath: c.filePath, actionFile: c.activeFile, fileSuffix: getFileSuffix()})
	data = viper.GetInt64(c.activeKey)
	return
}
