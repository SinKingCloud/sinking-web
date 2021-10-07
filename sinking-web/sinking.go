package sinking_web

import (
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
		engine      *Engine
	}

	ErrorHandel struct {
		NotFound func(*Context)
		Fail     func(c *Context, code int, message string)
	}

	Engine struct {
		*RouterGroup
		router             *router
		errorHandel        *ErrorHandel
		groups             []*RouterGroup
		htmlTemplates      *template.Template
		funcMap            template.FuncMap
		MaxMultipartMemory int64
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter(), MaxMultipartMemory: defaultMultipartMemory}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Recovery())
	if debug {
		engine.Use(Logger())
	}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) ANY(pattern string, handler HandlerFunc) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodTrace,
		http.MethodPatch,
		http.MethodHead,
	}
	for _, v := range methods {
		group.addRoute(v, pattern, handler)
	}
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute(http.MethodGet, pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute(http.MethodPost, pattern, handler)
}

func (group *RouterGroup) OPTIONS(pattern string, handler HandlerFunc) {
	group.addRoute(http.MethodOptions, pattern, handler)
}

func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.addRoute(http.MethodPut, pattern, handler)
}

func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute(http.MethodDelete, pattern, handler)
}

func (group *RouterGroup) HEAD(pattern string, handler HandlerFunc) {
	group.addRoute(http.MethodHead, pattern, handler)
}

func (group *RouterGroup) TRACE(pattern string, handler HandlerFunc) {
	group.addRoute(http.MethodTrace, pattern, handler)
}

func (group *RouterGroup) PATCH(pattern string, handler HandlerFunc) {
	group.addRoute(http.MethodPatch, pattern, handler)
}

func (group *RouterGroup) SetErrorHandle(handle *ErrorHandel) {
	group.engine.errorHandel = handle
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.SetStatus(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHtmlGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (group *RouterGroup) PROXY(pattern string, uri string, filter func(r *http.Request) *http.Request) {
	fun := func(c *Context) {
		Try(func() {
			target, err := url.Parse(uri)
			if err != nil {
				c.JSON(500, "url format error.")
				return
			}
			c.StatusCode = 200
			c.Request.Host = c.Request.URL.Host
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, pattern[0:strings.Index(pattern, "*")], "/", 1)
			filter(c.Request)
			proxy := httputil.NewSingleHostReverseProxy(target)
			proxy.ServeHTTP(c.Writer, c.Request)
		}, func(err interface{}) {
			c.StatusCode = 500
		})
	}
	group.ANY(pattern, fun)
}

func server(addr string, engine *Engine) *http.Server {
	Author(engine, addr)
	server := &http.Server{
		ReadTimeout:  readTimeOut,
		WriteTimeout: writeTimeout,
		Addr:         addr,
		Handler:      engine,
	}
	return server
}

func (engine *Engine) Run(addr string) (err error) {
	server := server(addr, engine)
	return server.ListenAndServe()
}

func (engine *Engine) RunTLS(addr string, certFile string, keyFile string) (err error) {
	server := server(addr, engine)
	return server.ListenAndServeTLS(certFile, keyFile)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
