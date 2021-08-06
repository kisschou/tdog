package tdog

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type (
	Response struct {
		context *Context
	}
)

func (r *Response) New(c *Context) *Response {
	r.context = c
	// c.BaseController.Res = r
	return r
}

func (r *Response) Success(data interface{}, msg string) {
	r.JSON(http.StatusOK, &H{
		"data":     data,
		"err_code": 0,
		"message":  msg,
	})
}

func (r *Response) Error(code string) {
	r.JSON(http.StatusInternalServerError, &H{
		"err_code": code,
	})
}

func (r *Response) JSON(code int, obj interface{}) {
	if code >= 0 {
		r.context.Writer.WriteHeader(code)
	}

	r.context.Writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(r.context.Writer)
	/*
		newObj := make(map[string]interface{})
		newObj = obj.(H)
		if code != http.StatusOK {
			if _, ok := newObj["code"]; ok {
				errorImpl := NewError(newObj["code"].(string))
				newObj["err_code"] = errorImpl.Code()
				newObj["message"] = errorImpl.Msg()

				// 附加消息返回
				if _, ok := newObj["_msg"]; ok {
					newObj["message"] = newObj["message"].(string) + newObj["_msg"].(string)
					delete(newObj, "_msg")
				}

				// 替换提示消息
				if _, ok := newObj["_replace"]; ok {
					newObj["message"] = newObj["_replace"].(string)
					delete(newObj, "_replace")
				}

				// 模版输出
				if _, ok := newObj["_template"]; ok {
					newObj["message"] = fmt.Sprintf(newObj["_template"].(string), newObj["message"].(string))
					delete(newObj, "_template")
				}

				delete(newObj, "code")
			}
		} else {
			newObj["err_code"] = 0
		}
		obj = newObj
	*/
	if err := encoder.Encode(obj); err != nil {
		r.context.Error(err, obj)
		http.Error(r.context.Writer, err.Error(), 500)
	}
}

func (r *Response) XML(code int, obj interface{}) {
	if code >= 0 {
		r.context.Writer.WriteHeader(code)
	}
	r.context.Writer.Header().Set("Content-Type", "application/xml")
	encoder := xml.NewEncoder(r.context.Writer)
	if err := encoder.Encode(obj); err != nil {
		r.context.Error(err, obj)
		http.Error(r.context.Writer, err.Error(), 500)
	}
}

func (r *Response) String(code int, msg string) {
	r.context.Writer.Header().Set("Content-Type", "text/plain")
	r.context.Writer.WriteHeader(code)
	r.context.Writer.Write([]byte(msg))
}

func (r *Response) Data(code int, data []byte) {
	r.context.Writer.WriteHeader(code)
	r.context.Writer.Write(data)
}

func (r *Response) Html(data interface{}) {
	path, _ := os.Getwd()
	t, _ := template.ParseFiles(path + "/app/views/" + r.context.Req.RequestURI + ".tpl")
	t.Execute(r.context.Writer, data)
}

func (r *Response) Captcha(code string) {
	d := make([]int, 4)
	for i := 0; i < len(code); i++ {
		iInt, _ := strconv.Atoi(code[i : i+1])
		d[i] = iInt
	}
	r.context.Writer.Header().Set("Content-Type", "image/png")
	NewImage(d, 100, 40).WriteTo(r.context.Writer)
}

func (r *Response) Redirect(uri string) {
	http.Redirect(r.context.Writer, r.context.Req, uri, http.StatusTemporaryRedirect)
}
