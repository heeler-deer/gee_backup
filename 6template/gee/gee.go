package gee

import (
	"log"
	"net/http"
	"path"
	"strings"
	"html/template"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
		engine      *Engine
	}
	Engine struct {
		router *router
		*RouterGroup
		groups []*RouterGroup
		htmlTemplates *template.Template
		funcMap template.FuncMap
	}
)

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	engine := r.engine
	newGroup := &RouterGroup{
		prefix: r.prefix + prefix,
		parent: r,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (r *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := r.prefix + comp
	r.engine.router.AddRouter(method, pattern, handler)
}

// GET defines the method to add GET request
func (r *RouterGroup) GET(pattern string, handler HandlerFunc) {
	log.Printf("%s", handler)
	r.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (r *RouterGroup) POST(pattern string, handler HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc

	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := NewContext(w, req)
	c.handlers = middlewares
	c.engine=engine
	engine.router.handle(c)
}

func (r *RouterGroup) Use(middlewares ...HandlerFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *RouterGroup) CreateStaticHandler(realtivePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(r.prefix, realtivePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context){
		file:=c.Param("filepath")
		if _,err:=fs.Open(file);err!=nil{
			c.Status((http.StatusNotFound))
			return
		}
		fileServer.ServeHTTP(c.Writer,c.Req)
	}
}




func (r *RouterGroup)Staic(relativePath string, root string){
	handler:=r.CreateStaticHandler(relativePath,http.Dir(root))
	urlPattern:=path.Join(relativePath,"/*filepath")
	r.GET(urlPattern,handler)
}






func (engine *Engine) SetFuncMap(funcmap template.FuncMap){
	engine.funcMap=funcmap
}


func(engine *Engine) LoadHtmlGlob(pattern string){
	engine.htmlTemplates=template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}


