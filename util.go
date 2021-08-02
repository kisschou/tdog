// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import (
	"bytes"
	"errors"
	"fmt"
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
type util struct {
}

// NewUtil init new util module.
func NewUtil() *util {
	return &util{}
}

// GetFilesBySuffix Gets the filename of all the specified suffixes from the specified path.
// given string filePath means file path of scan
// given string suffix means catch for same suffix
// returns []string files list of file name, file name has no suffix
// returns error err throw it if has errors
func (u *util) GetFilesBySuffix(filePath string, suffix string) (files []string, err error) {
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
func (u *util) FileExists(path string) bool {
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
func (u *util) IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile check path is file given string path returns true when it's
func (u *util) IsFile(path string) bool {
	return !u.IsDir(path)
}

// DirExistsAndCreate 检测目录是否存在, 不存在就创建
func (u *util) DirExistsAndCreate(path string) {
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
func (u *util) RandomStr(length int, randType ...int) string {
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
func (u *util) RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

// InArray 判断数组中是否存在某个值
// @param string dataType 标识needle的类型,可选:
//		[]string, []int, []int64, map[string]string, map[string]int, map[string]int64
// @param interface{} input 输入值,这个类型等于dataType的值类型相同
// @param interface{} needle 检索的集合
// @return bool
func (u *util) InArray(dataType string, input, needle interface{}) bool {
	defer Recover()

	switch dataType {
	case "[]string":
		for _, v := range needle.([]string) {
			if v == input.(string) {
				return true
			}
		}
		break

	case "[]int":
		for _, v := range needle.([]int) {
			if v == input.(int) {
				return true
			}
		}
		break

	case "[]int64":
		for _, v := range needle.([]int64) {
			if v == input.(int64) {
				return true
			}
		}
		break

	case "map[string]string":
		for _, v := range needle.(map[string]string) {
			if v == input.(string) {
				return true
			}
		}
		break

	case "map[string]int":
		for _, v := range needle.(map[string]int) {
			if v == input.(int) {
				return true
			}
		}
		break

	case "map[string]int64":
		for _, v := range needle.([]int64) {
			if v == input.(int64) {
				return true
			}
		}
		break
	}

	return false
}

// Isset 判断数组中是否存在某个键
// @param string dataType 标识needle的类型,可选:
//		[]string, []int, []int64, []interface{}
//		map[string]string, map[string]int, map[string]int64, map[string]interface{}
// @param interface{} input 输入值,这个类型等于dataType的键类型相同
// @param interface{} needle 检索的集合
// @return bool
func (u *util) Isset(dataType string, input, needle interface{}) bool {
	defer Recover()

	switch dataType {
	case "[]string":
		_ = needle.([]string)[input.(int)]
		break

	case "[]int":
		_ = needle.([]int)[input.(int)]
		break

	case "[]int64":
		_ = needle.([]int64)[input.(int)]
		break

	case "[]interface{}":
		_ = needle.([]interface{})[input.(int)]

	case "map[string]string":
		if _, ok := needle.(map[string]string)[input.(string)]; !ok {
			return false
		}
		break

	case "map[string]int":
		if _, ok := needle.(map[string]int)[input.(string)]; !ok {
			return false
		}
		break

	case "map[string]int64":
		if _, ok := needle.(map[string]int64)[input.(string)]; !ok {
			return false
		}
		break

	case "map[string]interface{}":
		if _, ok := needle.(map[string]interface{})[input.(string)]; !ok {
			return false
		}
		break
	}
	return true
}

// Empty 判断数组指定key的值是否为空, 数字则大于0
// @param string dataType 标识needle的类型,可选:
//		[]string, []int, []int64,
//		map[string]string, map[string]int, map[string]int64
// @param interface{} input 输入值,这个类型等于dataType的键类型相同
// @param interface{} needle 检索的集合
// @return bool
func (u *util) Empty(dataType string, input, needle interface{}) bool {
	if !u.Isset(dataType, input, needle) {
		return false
	}
	defer Recover()

	switch dataType {
	case "[]string":
		return len(needle.([]string)[input.(int)]) < 1

	case "[]int":
		return needle.([]int)[input.(int)] < 1

	case "[]int64":
		return needle.([]int64)[input.(int)] < 1

	case "map[string]string":
		return len(needle.(map[string]string)[input.(string)]) < 1

	case "map[string]int":
		return needle.(map[string]int)[input.(string)] < 1

	case "map[string]int64":
		return needle.(map[string]int64)[input.(string)] < 1
	}
	return false
}

// ArrayUnique 数组去掉重复键
// 提交什么类型过来就会返回什么类型，只不过要自己处理
// @param string dataType 标识input和返回值的类型,可选:
//		[]string, []int, []int64, []interface{}, map[string]string,
//		map[string]int, map[string]int64, map[string]interface{}
// @param interface{} input 需要处理的数组
// @return 处理完成的数组
func (u *util) ArrayUnique(dataType string, input interface{}) interface{} {
	defer Recover()

	switch dataType {
	case "[]string":
		res := make([]string, 0)
		for k, v := range input.([]string) {
			if !u.Isset(dataType, k, res) {
				res[k] = v
			}
		}
		return res

	case "[]int":
		res := make([]int, 0)
		for k, v := range input.([]int) {
			if !u.Isset(dataType, k, res) {
				res[k] = v
			}
		}
		return res

	case "[]int64":
		res := make([]int64, 0)
		for k, v := range input.([]int64) {
			if !u.Isset(dataType, k, res) {
				res[k] = v
			}
		}
		return res

	case "[]interface{}":
		res := make([]interface{}, 0)
		for k, v := range input.([]interface{}) {
			if !u.Isset(dataType, k, res) {
				res[k] = v
			}
		}
		return res

	case "map[string]string":
		res := make(map[string]string, 0)
		for k, v := range input.(map[string]string) {
			if !u.Isset(dataType, k, res) {
				res[k] = v
			}
		}
		return res

	case "map[string]int":
		res := make(map[string]int, 0)
		for k, v := range input.(map[string]int) {
			if !u.Isset(dataType, k, res) {
				res[k] = v
			}
		}
		return res

	case "map[string]int64":
		res := make(map[string]int64, 0)
		for k, v := range input.(map[string]int64) {
			if !u.Isset(dataType, k, res) {
				res[k] = v
			}
		}
		return res

	case "map[string]interface{}":
		res := make(map[string]interface{}, 0)
		for k, v := range input.(map[string]interface{}) {
			if !u.Isset(dataType, k, res) {
				res[k] = v
			}
		}
		return res
	}

	return nil
}

// ArrayMerge 数组合并
// 必须指定数据类型且所有数组必须同类型
// @param dataType string 传入数据的类型,可选:
//		[]string, []int, []interface{}, map[string]string, map[string]int, map[string]interface{}
// @param list ...interface{} 同类型的列表
// @return interface{} 返回的类型和dataType是一致的，
// 为nil表示错误，一般是传入了不同类型的数组造成的
func (u *util) ArrayMerge(dataType string, list ...interface{}) interface{} {
	if len(list) < 2 {
		return list[0]
	}

	defer Recover()

	switch dataType {
	case "[]string":
		res := make([]string, 0)
		for _, info := range list {
			for _, v := range info.([]string) {
				res = append(res, v)
			}
		}
		return res

	case "[]int":
		res := make([]int, 0)
		for _, info := range list {
			for _, v := range info.([]int) {
				res = append(res, v)
			}
		}
		return res

	case "[]interface{}":
		res := make([]interface{}, 0)
		for _, info := range list {
			for _, v := range info.([]interface{}) {
				res = append(res, v)
			}
		}
		return res

	case "map[string]string":
		res := make(map[string]string, 0)
		for _, info := range list {
			for k, v := range info.(map[string]string) {
				res[k] = v
			}
		}
		return res

	case "map[string]int":
		res := make(map[string]int, 0)
		for _, info := range list {
			for k, v := range info.(map[string]int) {
				res[k] = v
			}
		}
		return res

	case "map[string]interface{}":
		res := make(map[string]interface{}, 0)
		for _, info := range list {
			for k, v := range info.(map[string]interface{}) {
				res[k] = v
			}
		}
		return res
	}

	return nil
}

// Remove 切片删除指定index
// @param dataType string 传入数据的类型,可选:
//		[]string, []int, []int64, []interface{}
// @param slice interface{} 待处理的切片
// @param index int 准备去掉的index
// @return interface{} 返回的类型和dataType是一致的，
// 为nil表示错误，一般是传入了不同类型的数组造成的
func (u *util) Remove(dataType string, slice interface{}, index int) interface{} {
	defer Recover()
	switch dataType {
	case "[]string":
		res := make([]string, 0)
		for k, v := range slice.([]string) {
			if k != index {
				res[k] = v
			}
		}
		return res
	case "[]int":
		res := make([]int, 0)
		for k, v := range slice.([]int) {
			if k != index {
				res[k] = v
			}
		}
		return res
	case "[]int64":
		res := make([]int64, 0)
		for k, v := range slice.([]int64) {
			if k != index {
				res[k] = v
			}
		}
		return res
	case "[]interface{}":
		res := make([]interface{}, 0)
		for k, v := range slice.([]interface{}) {
			if k != index {
				res[k] = v
			}
		}
		return res
	}
	return []string{}
}

// VerifyEmail 验证邮件格式是否正确
// @return bool
func (u *util) VerifyEmail(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// VerifyPhone 验证手机格式是否正确
// @return bool
func (u *util) VerifyPhone(phone string) bool {
	pattern := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

// VerifyDate 校验日期的合理性(YYYY-MM-DD)
func (u *util) VerifyDate(input string) bool {
	pattern := `^([0-9]{4})-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(input)
}

// VerifyDateTime 校验日期时间的合理性(YYYY-MM-DD HH:mm:ss)
func (u *util) VerifyDateTime(input string) bool {
	pattern := `^([0-9]{4})-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01]) ([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(input)
}

// GetMachineId 获取设备id
// 通过网卡ipv4生成
func (u *util) GetMachineId() int64 {
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
func (u *util) StructToMap(obj interface{}) map[string]interface{} {
	defer Recover()
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
func (u *util) UrlSplit(url string) (protocol, domain string, port int) {
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

// UrlJoint url拼接
// @params string protocol 协议
// @params string domain 域名
// @params int port 端口
// @return string url 拼接完成的url
func (u *util) UrlJoint(protocol, domain string, port int) (url string) {
	url = protocol + "://" + domain
	if port != 80 {
		url += ":" + strconv.Itoa(port)
	}
	return
}

// SnakeString 驼峰转蛇形
// @params string s 需要转换的字符串
// @return string
func (u *util) SnakeString(s string) string {
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
func (u *util) CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := true // 首字母是否小写 // true: 小写; false: 大写
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
func (u *util) UcFirst(s string) string {
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

// checkPortAlived 检测端口是否已经占用
// @params int port 端口号
// @return bool
func (u *util) checkPortAlived(port int) bool {
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
func (u *util) GetPidByPort(port int) int {
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

// SetEnv 设置环境变量
// @params interface{} key
// 		当key为string时为键值
// 		当key为map[string]string时为设定的kv数据
// @params string value 设定值
// @return error 错误信息
func (u *util) SetEnv(key interface{}, value string) error {
	switch reflect.TypeOf(key).Kind().String() {
	case "string":
		return os.Setenv(key.(string), value)
	case "map":
		if reflect.TypeOf(key).String() == "map[string]string" {
			var err error
			for k, v := range key.(map[string]string) {
				err = os.Setenv(k, v)
			}
			return err
		}
		break
	}
	return errors.New("nil pointer")
}

// GetEnv 获取环境变量
func (u *util) GetEnv(keys ...string) map[string]string {
	result := make(map[string]string)
	for _, key := range keys {
		result[key] = os.Getenv(key)
	}
	return result
}

// Recover 从恐慌(panic)中走出来,并把造成恐慌的源头写入日志
// 这个函数一般用于defer
func Recover() {
	if err := recover(); err != nil {
		NewLogger().Error(fmt.Sprintf("%s", err))
	}
}

// Monitor 环境检测
// @return error err 错误信息
func Monitor() (err error) {
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
