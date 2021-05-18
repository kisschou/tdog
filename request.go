package tdog

import (
	"encoding/json"
	"mime/multipart"
	"net"
	"reflect"
	"strconv"
	"strings"
)

type (
	Request struct {
		Host    string
		IpAddr  string
		Header  map[string]string
		Params  map[string]string
		Get     map[string]string
		Post    map[string]string
		Put     map[string]interface{}
		File    *RequestFile
		IsGet   bool
		IsPost  bool
		IsPut   bool
		IsOpt   bool
		IsDel   bool
		IsPatch bool
	}

	RequestFile struct {
		Filename string
		Header   map[string][]string
		Size     int64
		Body     multipart.File
	}
)

func (r *Request) New(c *Context) *Request {
	var ghostMap map[string]string
	ghostMap = make(map[string]string)
	for k, v := range c.Req.Header {
		ghostMap[k] = v[0]
	}
	r.Header = ghostMap

	ghostMap = make(map[string]string)
	for k, v := range c.Req.URL.Query() {
		ghostMap[k] = v[0]
	}
	r.Get = ghostMap

	// 请求参数
	if _, ok := c.Req.Header["Content-Type"]; ok {
		if strings.Contains(c.Req.Header["Content-Type"][0], "json") {
			decoder := json.NewDecoder(c.Req.Body)
			var jsonParams map[string]interface{}
			_ = decoder.Decode(&jsonParams)
			r.Put = jsonParams
		}

		if strings.Contains(c.Req.Header["Content-Type"][0], "x-www-form-urlencoded") {
			for k, v := range c.Req.PostForm {
				ghostMap[k] = v[0]
			}
			r.Post = ghostMap
		}

		if strings.Contains(c.Req.Header["Content-Type"][0], "form-data") {
			fileBody, fileHeader, err := c.Req.FormFile("file")
			if err == nil {
				defer fileBody.Close()

				requestFile := new(RequestFile)
				requestFile.Filename = fileHeader.Filename
				requestFile.Header = fileHeader.Header
				requestFile.Size = fileHeader.Size
				requestFile.Body = fileBody
				r.File = requestFile
			}

			ghostMap = make(map[string]string)
			for k, v := range c.Req.PostForm {
				ghostMap[k] = v[0]
			}
			r.Post = ghostMap
		}
	}

	// Get|Put|Delete请求藏在地址中的参数
	if c.Req.Method == "GET" || c.Req.Method == "PUT" || c.Req.Method == "DELETE" {
		methodParams := make(map[string]string)
		for _, v := range c.Params {
			methodParams[v.Key] = v.Value
		}
		r.Get = NewUtil().ArrayMerge("map[string]string", r.Get, methodParams).(map[string]string)
	}

	// 获取请求地址
	r.Host = c.Req.Host

	// 获取客户端ip地址
	r = getIpAddr(r, c)

	// 判断请求类型
	r = checkReqMethod(r, c.Req.Method)

	// 合并参数到Params
	r = merge2Params(r)

	// set to base controller.
	// c.BaseController.Req = r
	return r
}

func checkReqMethod(r *Request, method string) *Request {
	switch method {
	case "GET":
		r.IsGet = true
		break
	case "POST":
		r.IsPost = true
		break
	case "PUT":
		r.IsPut = true
		break
	case "DELETE":
		r.IsDel = true
		break
	case "OPTIONS":
		r.IsOpt = true
		break
	case "PATCH":
		r.IsPatch = true
		break
	}
	return r
}

func merge2Params(r *Request) *Request {
	params := NewUtil().ArrayMerge("map[string]string", r.Get, r.Post).(map[string]string)
	if len(r.Put) > 0 {
		for k, v := range r.Put {
			if reflect.TypeOf(v).Kind().String() == "map" {
				dataJson, _ := json.Marshal(v.(map[string]interface{}))
				params[k] = string(dataJson)
			} else if reflect.TypeOf(v).Kind().String() == "int64" {
				params[k] = strconv.FormatInt(v.(int64), 10)
			} else if reflect.TypeOf(v).Kind().String() == "int" {
				params[k] = strconv.Itoa(v.(int))
			} else if reflect.TypeOf(v).Kind().String() == "float32" {
				params[k] = strconv.FormatFloat(v.(float64), 'f', -1, 32)
			} else if reflect.TypeOf(v).Kind().String() == "float64" {
				params[k] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
			} else if reflect.TypeOf(v).Kind().String() == "slice" {
				data := make([]map[string]interface{}, 0)
				for _, eachObj := range v.([]interface{}) {
					data = append(data, eachObj.(map[string]interface{}))
				}
				dataJson, _ := json.Marshal(data)
				params[k] = string(dataJson)
			} else {
				params[k] = v.(string)
			}
		}
	}
	r.Params = params
	return r
}

func getIpAddr(r *Request, c *Context) *Request {
	ip := ""

	ip = strings.TrimSpace(strings.Split(c.Req.Header.Get("X-Forwarded-For"), ",")[0])

	if ip == "" {
		ip = strings.TrimSpace(c.Req.Header.Get("X-Real-Ip"))
	}

	if ip == "" {
		var err error
		if ip, _, err = net.SplitHostPort(strings.TrimSpace(c.Req.RemoteAddr)); err != nil {
			ip = ""
		}
	}

	if ip == "::1" {
		ip = "127.0.0.1"
	}

	r.IpAddr = ip
	return r
}
