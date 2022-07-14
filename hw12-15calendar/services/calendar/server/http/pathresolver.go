package httpserver

import (
	"net/http"
)

type pathResolver struct {
	handlers map[string]func(res http.ResponseWriter, req *http.Request, server *Server)
}

func newPathResolver() *pathResolver {
	return &pathResolver{make(map[string]func(res http.ResponseWriter, req *http.Request, server *Server))}
}

func (p *pathResolver) Add(path string, handler func(res http.ResponseWriter, req *http.Request, server *Server)) {
	p.handlers[path] = handler
}
