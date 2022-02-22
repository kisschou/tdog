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
		Method string                 `json:"method"`
		Header map[string]string      `json:"header"`
		Url    string                 `json:"url"`
		Params map[string]interface{} `json:"params"`
	}

	httpResponse struct {
		Code        int                    `json:"http_code"`
		Data        string                 `json:"response"`
		Formatted   map[string]interface{} `json:"formatted"`
		ElapsedTime int64                  `json:"elapsed_time"`
	}
)

// NewRequest 初始化一个请求模型
func NewRequest() *httpRequest {
	return &httpRequest{Method: "", Header: map[string]string{}, Url: "", Params: map[string]interface{}{}}
}

func sendRequest(hp *httpRequest) (response *httpResponse, err error) {
	// 入参转json
	reqData, err := json.Marshal(hp.Params)
	if err != nil {
		return
	}

	// 创建请求体
	req, err := http.NewRequest(hp.Method, hp.Url, bytes.NewBuffer(reqData))
	if err != nil {
		return
	}

	// 写入请求头部
	for k, v := range hp.Header {
		req.Header.Set(k, v)
	}

	// 发起请求
	var res *http.Response
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	// 获取返回的内容
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	// 数据处理
	response = new(httpResponse)
	response.Code = res.StatusCode
	response.Data = string(body)
	formatted := make(map[string]interface{}, 0)
	_ = json.Unmarshal(body, &formatted)
	response.Formatted = formatted
	return
}

// Send 发起表单请求
func (hp *httpRequest) Send() *httpResponse {
	startTime := time.Now().UnixNano()
	response, err := sendRequest(hp)
	info, err := json.Marshal(&httpRequest{
		Method: hp.Method,
		Header: hp.Header,
		Url:    hp.Url,
		Params: hp.Params,
	})
	logInfo := "请求: " + string(info)
	if err != nil {
		// 写入log
		logInfo += "时发生错误, 错误信息: " + err.Error()
		go NewLogger().Error(logInfo)

		errorImpl := NewError("SERVICE_ERROR")
		var formatted = map[string]interface{}{
			"err_code": errorImpl.Code(),
			"message":  errorImpl.Msg(),
		}
		data, _ := json.Marshal(formatted)
		return &httpResponse{
			Code:        http.StatusInternalServerError,
			Formatted:   formatted,
			Data:        string(data),
			ElapsedTime: time.Now().UnixNano() - startTime,
		}
	}
	response.ElapsedTime = time.Now().UnixNano() - startTime

	// 写入log
	// data, _ := json.Marshal(response)
	// logInfo += "成功， 返回数据: " + string(data)
	// go NewLogger().Info(logInfo)

	return response
}
