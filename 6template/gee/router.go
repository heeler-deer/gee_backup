package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
	roots map[string ]*node
}

func NewRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc),roots:make(map[string]*node)}
}

//只能允许一个实例
func ParsePattern(pattern string) []string{
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}




func (r *router) AddRouter(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := ParsePattern(pattern)
	key := method + "-" + pattern
	_,ok:=r.roots[method]
	if !ok{
		r.roots[method]=&node{}
	}
	r.roots[method].Insert(pattern,parts,0)
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	n, params := r.GetRoute(c.Method, c.Path)

	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern

		c.handlers=append(c.handlers,r.handlers[key])
		log.Printf("%s",c.handlers)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}

func (r *router) GetRoute(method string, path string) (*node, map[string]string){
	searchParts:=ParsePattern(path)
	params:=make(map[string]string)
	root, ok:=r.roots[method]
	if !ok{
		return nil,nil
	}
	n:=root.Search(searchParts,0)

	if n!=nil{
		parts:=ParsePattern(n.pattern)
		for index,part:=range parts{
			log.Printf("parts: %s , searchParts: %s",parts,searchParts)
			if part[0]==':'{
				params[part[1:]]=searchParts[index]

			}
			if part[0]=='*' && len(part)>1{
				params[part[1:]] = strings.Join(searchParts[index:],"/")
				break
			}
		}
		return n, params
	}
	return nil,nil
} 