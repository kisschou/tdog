package tdog

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	httpRequest struct {
		Method string
		Header map[string]string
		Url    string
		Params map[string]interface{}
	}

	httpResponse struct {
		Code        int
		Data        string
		Formatted   map[string]interface{}
		ElapsedTime int64
		Error       error
	}
)

// NewRequest 初始化一个请求模型
func NewRequest() *httpRequest {
	return &httpRequest{Method: "", Header: map[string]string{}, Url: "", Params: map[string]interface{}{}}
}

// FormRequest 发起表单请求
func (hp *httpRequest) FormRequest() (httpCode int, resData string, elapsedTime int64, err error) {
	startTime := time.Now().UnixNano()
	client := &http.Client{}
	reqDataJson, _ := json.Marshal(hp.Params)
	req, err := http.NewRequest(hp.Method, hp.Url, bytes.NewBuffer(reqDataJson))
	if err != nil {
		httpCode = http.StatusInternalServerError
		elapsedTime = time.Now().UnixNano() - startTime
		return
	}
	for k, v := range hp.Header {
		req.Header.Set(k, v)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		httpCode = http.StatusInternalServerError
		elapsedTime = time.Now().UnixNano() - startTime
		return
	}
	elapsedTime = time.Now().UnixNano() - startTime
	httpCode = res.StatusCode
	resData = string(body)
	return
}

// BytesPost 发送二进制数据流
func (hp *httpRequest) BytesPost() (int, string, int64, error) {
	startTime := time.Now().UnixNano()
	var elapsedTime int64
	data, _ := json.Marshal(hp.Params)
	body := bytes.NewReader(data)
	req, err := http.NewRequest(hp.Method, hp.Url, body)
	if err != nil {
		go NewLogger().Error(err.Error())
		elapsedTime = time.Now().UnixNano() - startTime
		return http.StatusInternalServerError, "", elapsedTime, err
	}

	for k, v := range hp.Header {
		req.Header.Set(k, v)
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		go NewLogger().Error(err.Error())
		elapsedTime = time.Now().UnixNano() - startTime
		return http.StatusInternalServerError, "", elapsedTime, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		go NewLogger().Error(err.Error())
	}

	elapsedTime = time.Now().UnixNano() - startTime
	return resp.StatusCode, string(b), elapsedTime, err
}

// ServicePost 针对请求网关服务构建
// 继承后只需要set结构体中的Params
func (hp *httpRequest) ServicePost() (bool, string, int64) {
	// header
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Connection"] = "Keep-Alive"

	hp.Method = "POST"
	hp.Header = header
	hp.Url = NewConfig().Get("api_url.gateway_url").ToString() + "/feign/http"
	httpCode, res, elapsedTime, err := hp.BytesPost()
	if httpCode != http.StatusOK || err != nil {
		if err != nil {
			go NewLogger().Error(err.Error())
		}
		return false, res, elapsedTime
	}
	return true, res, elapsedTime
}
