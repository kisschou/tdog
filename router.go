package tdog

import (
	"log"
	"math"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/kisschou/ReqHelper"
)

const (
	// AbortIndex .
	AbortIndex = math.MaxInt8 / 2
)

type (
	// HandlerFunc .
	HandlerFunc func(*HttpUtil)

	// H .
	// H map[string]interface{}

	// ErrorMsg .
	ErrorMsg struct {
		Message string      `json:"msg"`
		Meta    interface{} `json:"meta"`
	}

	// Context .
	Context struct {
		Req      *http.Request
		Writer   http.ResponseWriter
		Keys     map[string]interface{}
		Errors   []ErrorMsg
		Params   httprouter.Params
		handlers []HandlerFunc
		engine   *HttpEngine
		index    int8
	}

	// RouterGroup .
	RouterGroup struct {
		Handlers []HandlerFunc
		prefix   string
		parent   *RouterGroup
		engine   *HttpEngine
	}

	// HttpEngine .
	HttpEngine struct {
		*RouterGroup
		router *httprouter.Router
	}

	HttpUtil struct {
		Req *ReqHelper.Request
		Res *Response
	}

	// H ResponseText
	H map[string]interface{}
)

// New HttpEngine
func New() *HttpEngine {
	engine := &HttpEngine{}
	engine.RouterGroup = &RouterGroup{nil, "", nil, engine}
	engine.router = httprouter.New()

	// 静态资源
	engine.router.ServeFiles("/data/upload/*filepath", http.Dir("data/upload/"))

	return engine
}

// Default Returns a Engine instance with the Logger and Recovery already attached.
func Default() *HttpEngine {
	engine := New()
	// engine.Use(Recovery(), Logger())
	return engine
}

// ServeFiles  router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (engine *HttpEngine) ServeFiles(path string, root http.FileSystem) {
	engine.router.ServeFiles(path, root)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (engine *HttpEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL)
	// 跨域问题
	// 如果请求为OPTIONS(跨域请求)返回跨域授权
	if req.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		return
	}
	engine.router.ServeHTTP(w, req)
}

// Run .
func (engine *HttpEngine) Run() {
	Println("Start Listen :"+NewConfig().Get("app_port").ToString(), 41)
	http.ListenAndServe(":"+NewConfig().Get("app_port").ToString(), engine)
}

/************************************/
/********** ROUTES GROUPING *********/
/************************************/

func (group *RouterGroup) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handlers []HandlerFunc) *Context {
	return &Context{
		Writer:   w,
		Req:      req,
		index:    -1,
		engine:   group.engine,
		Params:   params,
		handlers: handlers,
	}
}

// Use .
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middlewares...)
}

// Group .
func (group *RouterGroup) Group(component string, handlers ...HandlerFunc) *RouterGroup {
	prefix := path.Join(group.prefix, component)
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		parent:   group,
		prefix:   prefix,
		engine:   group.engine,
	}
}

// Handle .
func (group *RouterGroup) Handle(method, p string, handlers []HandlerFunc) {
	p = path.Join(group.prefix, p)
	handlers = group.combineHandlers(handlers)
	group.engine.router.Handle(method, p, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		group.createContext(w, req, params, handlers).Next()
	})
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (group *RouterGroup) POST(path string, handlers ...HandlerFunc) {
	group.Handle("POST", path, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (group *RouterGroup) GET(path string, handlers ...HandlerFunc) {
	group.Handle("GET", path, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (group *RouterGroup) DELETE(path string, handlers ...HandlerFunc) {
	group.Handle("DELETE", path, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (group *RouterGroup) PATCH(path string, handlers ...HandlerFunc) {
	group.Handle("PATCH", path, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (group *RouterGroup) PUT(path string, handlers ...HandlerFunc) {
	group.Handle("PUT", path, handlers)
}

func (group *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	s := len(group.Handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	h = append(h, group.Handlers...)
	h = append(h, handlers...)
	return h
}

// Next .
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		req := ReqHelper.New(c.Req)
		for _, v := range c.Params {
			req.Extra[v.Key] = v.Value
		}
		c.handlers[c.index](&HttpUtil{
			Req: req,
			Res: new(Response).New(c),
		})
	}
}

// Abort .
func (c *Context) Abort(code int) {
	c.Writer.WriteHeader(code)
	c.index = AbortIndex
}

// Fail .
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

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set Sets a new pair key/value just for the specefied context.
// It also lazy initializes the hashmap
func (c *Context) Set(key string, item interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = item
}

// Get Returns the value for the given key.
// It panics if the value doesn't exist.
func (c *Context) Get(key string) interface{} {
	var ok bool
	var item interface{}
	if c.Keys != nil {
		item, ok = c.Keys[key]
	} else {
		item, ok = nil, false
	}
	if !ok || item == nil {
		log.Panicf("Key %s doesn't exist", key)
	}
	return item
}
