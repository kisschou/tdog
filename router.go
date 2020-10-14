package tdog

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type (
	// HandlerFunc .
	HandlerFunc func()

	// ErrorMsg
	ErrorMsg struct {
		Message string      `json:"msg"`
		Meta    interface{} `json:"meta"`
	}

	// Context .
	Context struct {
		Req            *http.Request
		Writer         http.ResponseWriter
		Errors         []ErrorMsg
		Params         httprouter.Params
		handler        HandlerFunc
		engine         *HttpEngine
		BaseController *Controller
	}

	// RouterGroup .
	RouterGroup struct {
		Handler        HandlerFunc
		prefix         string
		parent         *RouterGroup
		engine         *HttpEngine
		BaseController *Controller
	}

	// HttpEngine .
	HttpEngine struct {
		*RouterGroup
		router *httprouter.Router
	}
)

// New HttpEngine
func NewEngine() *HttpEngine {
	engine := &HttpEngine{}
	engine.RouterGroup = &RouterGroup{nil, "", nil, engine, nil}
	engine.router = httprouter.New()

	// 静态资源
	engine.router.ServeFiles("/public/static/*filepath", http.Dir("public/static/"))

	return engine
}

// ServeHTTP makes the router implement the http.Handler interface.
func (engine *HttpEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// OPTIONS请求直接返回
	if req.Method == "OPTIONS" {
		// 跨域
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		return
	}

	engine.router.ServeHTTP(w, req)
}

// Run .
func (engine *HttpEngine) Run() {
	ConfigLib := new(Config)
	http.ListenAndServe(":"+ConfigLib.Get("app_port").String(), engine)
}

/************************************/
/********** ROUTES GROUPING *********/
/************************************/

func (group *RouterGroup) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handler HandlerFunc) *Context {
	if _, ok := req.Header["Content-Type"]; ok {
		if strings.Contains(req.Header["Content-Type"][0], "x-www-form-urlencoded") {
			req.ParseForm()
		}

		if strings.Contains(req.Header["Content-Type"][0], "form-data") {
			req.ParseMultipartForm(32 << 20)
		}
	}

	// 跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	return &Context{
		Writer:         w,
		Req:            req,
		engine:         group.engine,
		Params:         params,
		handler:        handler,
		BaseController: group.BaseController,
	}
}

// Group .
func (group *RouterGroup) Group(component string, baseController *Controller) *RouterGroup {
	prefix := path.Join(group.prefix, component)
	return &RouterGroup{
		Handler:        nil,
		parent:         group,
		prefix:         prefix,
		engine:         group.engine,
		BaseController: baseController,
	}
}

// Handle .
func (group *RouterGroup) Handle(method, p string, handler HandlerFunc) {
	p = path.Join(group.prefix, p)
	group.engine.router.Handle(method, p, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		group.createContext(w, req, params, handler).Next()
	})
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (group *RouterGroup) POST(path string, handler HandlerFunc, baseController ...*Controller) {
	if len(baseController) > 0 {
		group.BaseController = baseController[0]
	}
	group.Handle("POST", path, handler)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (group *RouterGroup) GET(path string, handler HandlerFunc, baseController ...*Controller) {
	if len(baseController) > 0 {
		group.BaseController = baseController[0]
	}
	group.Handle("GET", path, handler)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (group *RouterGroup) DELETE(path string, handler HandlerFunc, baseController ...*Controller) {
	if len(baseController) > 0 {
		group.BaseController = baseController[0]
	}
	group.Handle("DELETE", path, handler)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (group *RouterGroup) PATCH(path string, handler HandlerFunc, baseController ...*Controller) {
	if len(baseController) > 0 {
		group.BaseController = baseController[0]
	}
	group.Handle("PATCH", path, handler)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (group *RouterGroup) PUT(path string, handler HandlerFunc, baseController ...*Controller) {
	if len(baseController) > 0 {
		group.BaseController = baseController[0]
	}
	group.Handle("PUT", path, handler)
}

// OPTIONS
func (group *RouterGroup) OPTIONS(path string, handler HandlerFunc, baseController ...*Controller) {
	if len(baseController) > 0 {
		group.BaseController = baseController[0]
	}
	group.Handle("OPTIONS", path, handler)
}

// Next .
func (c *Context) Next() {
	req := new(Request)
	req.New(c)
	res := new(Response)
	res.New(c)

	/**
	 * 验签
	 * 用户请求接口时候自动判断执行验签
	 **/
	ConfigLib := new(Config)
	// 跳过验签路由列表
	ignoreRouter := ConfigLib.Get("ignore_auth").StringSlice()
	isAuth := ConfigLib.Get("is_auth").Bool()
	for _, v := range ignoreRouter {
		if v == c.Req.RequestURI {
			isAuth = false
			break
		}
	}
	// 跳过登入校验的路由列表
	ignoreLoginRouter := ConfigLib.Get("ignore_login").StringSlice()
	isLogin := ConfigLib.Get("is_login").Bool()
	for _, v := range ignoreLoginRouter {
		if v == c.Req.RequestURI {
			isLogin = false
			break
		}
	}

	// 开始鉴权
	if isAuth {
		authorization := ""
		if _, ok := c.BaseController.Req.Header["Authorization"]; ok {
			if len(c.BaseController.Req.Header["Authorization"]) > 0 {
				authorization = c.BaseController.Req.Header["Authorization"][0]
			}
		}

		if len(authorization) < 1 {
			c.BaseController.Res.JSON(http.StatusUnauthorized, H{
				"code": "ERROR_UNAUTHOZED",
			})
			return
		}

		JwtCore := new(Jwt)
		if !JwtCore.Check(authorization) {
			c.BaseController.Res.JSON(http.StatusUnauthorized, H{
				"code": "ERROR_UNAUTHOZED",
			})
			return
		}

		// 开始登录校验
		if isLogin {
			var userId int64
			var err error
			if JwtCore.Get(authorization, "user_id") != nil {
				userId, err = strconv.ParseInt(JwtCore.Get(authorization, "user_id").(string), 10, 64)
				if err != nil {
					c.BaseController.Res.JSON(http.StatusInternalServerError, H{
						"code": "ERROR_UNLOGIN",
					})
					return
				}
			} else {
				c.BaseController.Res.JSON(http.StatusInternalServerError, H{
					"code": "ERROR_UNLOGIN",
				})
				return
			}
			c.BaseController.UserId = userId
		}
	}

	c.handler()
}

func (c *Context) Abort(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) Fail(code int, err error) {
	c.Error(err, "Operation aborted")
	c.Abort(code)
}

func (c *Context) Error(err error, meta interface{}) {
	c.Errors = append(c.Errors, ErrorMsg{
		Message: err.Error(),
		Meta:    meta,
	})
}
