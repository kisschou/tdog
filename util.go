package tdog

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Util struct {
}

// 判断文件是否存在
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

// 判断是否是目录
func (u *Util) IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断是否是文件
func (u *Util) IsFile(path string) bool {
	return !u.IsDir(path)
}

// 判断文件夹是否存在,不存在就创建
func (u *Util) DirExistsAndCreate(path string) {
	if !u.FileExists(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			logger := Logger{Level: 0, Key: "error"}
			logger.New(err.Error())
			os.Exit(0)
		}
	}
}

// 生成指定数量随机字母加数字
func (u *Util) RandomStr(length int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 生成指定数量随机数字
func (u *Util) RandomNum(length int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func codeToKey(hashCode string) []byte {
	cryptMap := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(cryptMap)
	result := []byte{}
	var i int
	startIndex := 0
	endIndex := startIndex + 2
	for i = 0; i < (len(hashCode)/2 + len(hashCode)%2); i++ {
		if len(hashCode[startIndex:]) < 2 {
			si, _ := strconv.Atoi(hashCode[startIndex:])
			result = append(result, bytes[si])
			break
		}
		s := hashCode[startIndex:endIndex]
		si, _ := strconv.Atoi(s)
		if si > 61 {
			si, _ = strconv.Atoi(hashCode[startIndex : startIndex+1])
			result = append(result, bytes[si])
			res := codeToKey(hashCode[startIndex+1:])
			result = append(result, res...)
			break
		} else {
			result = append(result, bytes[si])
		}
		startIndex += 2
		endIndex = startIndex + 2
	}
	return result
}

func (u *Util) ShorturlKey(baseUrl string) string {
	CryptLib := new(Crypt)
	CryptLib.Str = baseUrl
	hashCode := CryptLib.Crc32()
	return string(codeToKey(hashCode))
}

// 判断map中是否存在某个key
func (u *Util) InMap(key string, dataMap map[string]interface{}) bool {
	if _, ok := dataMap[key]; ok {
		return true
	}
	return false
}

func (u *Util) InMapStringSlice(key string, dataMap map[string][]string) bool {
	if _, ok := dataMap[key]; ok {
		if len(dataMap[key]) > 0 {
			return true
		}
	}
	return false
}

func (u *Util) InStringSlice(key string, dataStringSlice []string) bool {
	for _, v := range dataStringSlice {
		if v == key {
			return true
		}
	}
	return false
}

// 检测字符串是邮件、手机号、字符串
// 返回: 0字符串1邮件2手机号
func (u *Util) CheckStrType(str string) int {
	if verifyEmailFormat(str) {
		return 1
	}

	if verifyPhoneFormat(str) {
		return 2
	}

	return 0
}

func verifyEmailFormat(email string) bool {
	// pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func verifyPhoneFormat(phone string) bool {
	pattern := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

// 获取设备id
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

func (u *Util) RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

func (u *Util) StructToMap(obj interface{}) map[string]interface{} {
	mapVal := make(map[string]interface{})
	elem := reflect.ValueOf(obj).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		mapVal[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return mapVal
}

// url分解
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

// url整合
func (u *Util) UrlJoint(protocol, domain string, port int) (url string) {
	url = protocol + "://" + domain
	if port != 80 {
		url += ":" + strconv.Itoa(port)
	}
	return
}

// 环境检测
func (u *Util) Monitor() (err error) {
	return
	// MySQL环境
	MySqlTdog := new(MySql)
	MySqlTdog.NewEngine()
	if MySqlTdog.Engine == nil || !MySqlTdog.Ping() {
		err = errors.New("ERROR: MySQL connect fail! Please start mysql server and retry!")
		return
	}
	// Redis环境
	RedisTdog := new(Redis)
	if RedisTdog.Engine == nil || !RedisTdog.Ping() {
		err = errors.New("ERROR: Redis connect fail! Please start redis server and retry!")
		return
	}
	return
}

// 获取雪花id
func (u *Util) GetSnowFlake() int64 {
	SnowflakeTdog, _ := NewSnowFlake(u.RandInt64(1, 1023))
	return SnowflakeTdog.GetId()
}

func (u *Util) Request(authorization, apiCode string, params map[string]interface{}) (data map[string]interface{}, err error) {
	FeignTdog := new(Feign)
	header := make(map[string]string)
	header["Authorization"] = authorization
	if len(params) > 0 {
		header["Content-Type"] = "application/json"
	}
	reqParams := make(map[string]interface{})
	reqParams["api_code"] = apiCode
	reqParams["params"] = params
	code, res, _ := FeignTdog.Url("http://127.0.0.1:8001/gateway/feign").Method("POST").Header(header).Params(reqParams).Target()
	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		return
	}
	if code != 200 {
		err = errors.New(data["message"].(string))
		return
	}
	return
}

func (u *Util) GetUserId(authorization string) (userId int64, err error) {
	ConfigTdog := new(Config)
	apiPath := ConfigTdog.Get("api_path").String()
	FeignTdog := new(Feign)
	header := make(map[string]string)
	header["Authorization"] = authorization
	code, res, _ := FeignTdog.Url(apiPath + ":8001/gateway/auth/getKey/user_id").Method("GET").Header(header).Target()
	if code == http.StatusOK {
		dataMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(res), &dataMap)
		if err != nil {
			err = errors.New("ERROR_UNLOGIN")
			return
		}
		switch dataMap["value"].(type) {
		case float64:
			userId = int64(dataMap["value"].(float64))
		case string:
			userId, err = strconv.ParseInt(dataMap["value"].(string), 10, 64)
			if err != nil {
				err = errors.New("ERROR_UNLOGIN")
				return
			}
		}
	} else {
		err = errors.New("ERROR_UNLOGIN")
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

/**
 * 驼峰转蛇形 snake string
 * @param s 需要转换的字符串
 * @return string
 **/
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

/**
 * 蛇形转驼峰
 * @param s要转换的字符串
 * @return string
 */
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

/**
 * 首字母转大写
 * @param s要转换的字符串
 * @return string
 **/
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
