package tdog

type (
	Feign struct {
		method string
		url    string
		header map[string]string
		params map[string]interface{}
	}
)

func NewFeign() *Feign {
	feign := &Feign{"GET", "", nil, nil}
	return feign
}

func (feign *Feign) Url(url string) *Feign {
	feign.url = url
	return feign
}

func (feign *Feign) Header(header map[string]string) *Feign {
	feign.header = header
	return feign
}

func (feign *Feign) Method(method string) *Feign {
	feign.method = method
	return feign
}

func (feign *Feign) Params(params map[string]interface{}) *Feign {
	feign.params = params
	return feign
}

func (feign *Feign) Target() (code int, res string, elapsedTime int64) {
	var err error
	HttpRequestLib := new(HttpRequest)
	HttpRequestLib.Method = feign.method
	HttpRequestLib.Header = feign.header
	HttpRequestLib.Url = feign.url
	HttpRequestLib.Params = feign.params
	code, res, elapsedTime, err = HttpRequestLib.FormRequest()
	if err != nil {
		code = 0
		res = "ERROR_FEIGN_REQUEST_FAIL"

		go NewLogger().Error(err.Error())
		return
	}
	return
}
