// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/**
 * The module for uril handler.
 *
 * Author: Kisschou
 * @Build: 2020-11-22
 */
type Util struct {
}

// NewUtil init new util module.
func NewUtil() *Util {
	return &Util{}
}

// get list of file's name from path with suffix.
// given string filePath means file path of scan
// given string suffix means catch for same suffix
// returns []string files list of file name, file name has no suffix
// returns error err throw it if has errors
func (u *Util) GetFilesBySuffix(filePath string, suffix string) (files []string, err error) {
	rd, err := ioutil.ReadDir(filePath)
	if err != nil {
		NewLogger().Error(err.Error())
		return
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			fileSuffix := path.Ext(fi.Name())
			fileName := strings.TrimSuffix(fi.Name(), fileSuffix)
			if "."+suffix == fileSuffix {
				files = append(files, fileName)
			}
		}
	}
	return
}

// FileExists check file is exists. given string path returns true when exists
func (u *Util) FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir check path is dir given string path returns true when it's
func (u *Util) IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile check path is file given string path returns true when it's
func (u *Util) IsFile(path string) bool {
	return !u.IsDir(path)
}

// DirExistsAndCreate 检测目录是否存在, 不存在就创建
func (u *Util) DirExistsAndCreate(path string) {
	if !u.FileExists(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			NewLogger().Error(err.Error())
			os.Exit(0)
		}
	}
}

// RandomStr 生成指定长度的字符串
// @params int length 长度
// @params ...int randType 默认值为纯数字1加入数字2加入小写字母3加入大写字母
// @return string 结果
func (u *Util) RandomStr(length int, randType ...int) string {
	str := ""
	for _, v := range randType {
		switch v {
		case 1:
			str += "0123456789"
			break
		case 2:
			str += "abcdefghijklmnopqrstuvwxyz"
			break
		case 3:
			str += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
			break
		}
	}
	if len(str) < 8 {
		str = "0123456789"
	}

	bytes := []byte(str)
	result := []byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// RandInt64 生成指定范围的随机int64类型数据
// @params int64 min 最小值
// @params int64 max 最大值
// @return int64 结果
func (u *Util) RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

// InMap 判断map中是否存在某个key
// @params string key 键名
// @params map[string]interface{} needle
// @return bool true when exists
func (u *Util) InMap(key string, dataMap map[string]interface{}) bool {
	if _, ok := dataMap[key]; ok {
		return true
	}
	return false
}

// InMapStringSlice same of InMap, just diff style of needle
func (u *Util) InMapStringSlice(key string, dataMap map[string][]string) bool {
	if _, ok := dataMap[key]; ok {
		if len(dataMap[key]) > 0 {
			return true
		}
	}
	return false
}

// InStringSlice same of InMap, just diff style of needle
func (u *Util) InStringSlice(key string, dataStringSlice []string) bool {
	for _, v := range dataStringSlice {
		if v == key {
			return true
		}
	}
	return false
}

// CheckStrType 检测字符串是邮件、手机号、字符串
// @return 0字符串1邮件2手机号
func (u *Util) CheckStrType(str string) int {
	if u.VerifyEmail(str) {
		return 1
	}

	if u.VerifyPhone(str) {
		return 2
	}

	return 0
}

// VerifyEmailFormat 验证邮件格式是否正确
// @return bool
func (u *Util) VerifyEmail(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// VerifyPhoneFormat 验证手机格式是否正确
// @return bool
func (u *Util) VerifyPhone(phone string) bool {
	pattern := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

// GetMachineId 获取设备id
// 通过网卡ipv4生成
func (u *Util) GetMachineId() int64 {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 127000000001
	}

	gatewayIp := ""
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				gatewayIp = ipnet.IP.String()
				break
			}
		}
	}

	gatewayIpSplit := strings.Split(gatewayIp, ".")
	gatewayIp = ""
	for _, v := range gatewayIpSplit {
		if len(v) < 3 {
			for i := len(v); i < 3; i++ {
				v = "0" + v
			}
		}
		gatewayIp += v
	}
	machineId, _ := strconv.ParseInt(gatewayIp, 10, 64)
	return machineId
}

// StructToMap 结构体转换成map
func (u *Util) StructToMap(obj interface{}) map[string]interface{} {
	mapVal := make(map[string]interface{})
	elem := reflect.ValueOf(obj).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		mapVal[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return mapVal
}

// UrlSplit url分解
// @params string url 需要拆解的url
// @return string protocol 协议
// @return string domain 域名
// @return int port 端口
func (u *Util) UrlSplit(url string) (protocol, domain string, port int) {
	var err error
	urlCompose := strings.Split(url, ":")
	protocol = urlCompose[0]
	domain = strings.Replace(urlCompose[1], "/", "", -1)
	port = 80
	if len(urlCompose) > 2 {
		port, err = strconv.Atoi(strings.Replace(urlCompose[2], "/", "", -1))
		if err != nil {
			port = 80
		}
	}
	return
}

// UrlJoint url整合
// @params string protocol 协议
// @params string domain 域名
// @params int port 端口
// @return string url 需要拆解的url
func (u *Util) UrlJoint(protocol, domain string, port int) (url string) {
	url = protocol + "://" + domain
	if port != 80 {
		url += ":" + strconv.Itoa(port)
	}
	return
}

// Monitor 环境检测
func (u *Util) Monitor() (err error) {
	return
	// MySQL环境
	if NewMySQL().Engine == nil {
		err = errors.New("ERROR: MySQL connect fail! Please start mysql server and retry!")
		NewLogger().Error(err.Error())
		return
	}
	if err = NewMySQL().Engine.Ping(); err != nil {
		NewLogger().Error(err.Error())
		err = errors.New("ERROR: MySQL connect fail! Please start mysql server and retry!")
		return
	}
	// Redis环境
	if NewRedis().Engine == nil {
		err = errors.New("ERROR: Redis connect fail! Please start redis server and retry!")
		return
	}
	return
}

func (u *Util) MySQLColumnTypeConvert(columnType string) string {
	convert := "string"
	switch strings.ToUpper(columnType) {
	case "CHAR", "VARCHAR":
		convert = "string"
		break

	case "TINYBLOB", "TINYTEXT", "BLOB", "TEXT", "MEDIUMBLOB", "MEDIUMTEXT", "LONGBLOB", "LONGTEXT":
		convert = "text"
		break

	case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT":
		convert = "select"
		break

	case "FLOAT", "DOUBLE":
		convert = "string"
		break

	case "DATE", "TIME", "YEAR", "DATETIME", "TIMESTAMP":
		convert = "date"
		break
	}
	return convert
}

// SnakeString 驼峰转蛇形
// @params string s 需要转换的字符串
// @return string
func (u *Util) SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

// CamelString 蛇形转驼峰
// @param s 要转换的字符串
// @return string
func (u *Util) CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// UcFirst 首字母转大写
// @param s要转换的字符串
// @return string
func (u *Util) UcFirst(s string) string {
	var upperStr string
	vv := []rune(s)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				return s
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

// checkPortAlived 检测端口是否已经暂用
// @params int port 端口号
// @return bool
func (u *Util) checkPortAlived(port int) bool {
	isAlived := false
	var outBytes bytes.Buffer
	cmd0 := exec.Command("netstat", "-ano")
	cmd1 := exec.Command("grep", strconv.Itoa(port))
	cmd1.Stdin, _ = cmd0.StdoutPipe()
	cmd2 := exec.Command("grep", "LISTEN")
	cmd2.Stdin, _ = cmd1.StdoutPipe()
	cmd2.Stdout = &outBytes
	_ = cmd2.Start()
	_ = cmd1.Start()
	_ = cmd0.Run()
	_ = cmd1.Wait()
	_ = cmd2.Wait()
	res := outBytes.String()
	if len(res) > 10 {
		isAlived = true
	}
	return isAlived
}

// GetPidByPort 通过端口号获取pid
// 不是一定能获取到，如协程就获取不到
// @params int port 端口号
// @return int pid
func (u *Util) GetPidByPort(port int) int {
	pid := -1
	var outBytes bytes.Buffer
	cmd0 := exec.Command("netstat", "-ano")
	cmd1 := exec.Command("grep", strconv.Itoa(port))
	cmd1.Stdin, _ = cmd0.StdoutPipe()
	cmd2 := exec.Command("grep", "LISTENING")
	cmd2.Stdin, _ = cmd1.StdoutPipe()
	cmd2.Stdout = &outBytes
	_ = cmd2.Start()
	_ = cmd1.Start()
	_ = cmd0.Run()
	_ = cmd1.Wait()
	_ = cmd2.Wait()
	res := outBytes.String()
	r := regexp.MustCompile(`\s\d+\s`).FindAllString(res, -1)
	if len(r) > 0 {
		var err error
		pid, err = strconv.Atoi(strings.TrimSpace(r[1]))
		if err != nil {
			pid = -1
		}
	}
	return pid
}
