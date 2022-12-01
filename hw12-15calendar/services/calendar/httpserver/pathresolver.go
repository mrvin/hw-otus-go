package httpserver

import (
	"fmt"
	"net/http"
	"path"
)

// TODO:add thread safety.
type pathResolver struct {
	handlers map[string]func(res http.ResponseWriter, req *http.Request, server *Server)
}

func newPathResolver() *pathResolver {
	return &pathResolver{make(map[string]func(res http.ResponseWriter, req *http.Request, server *Server))}
}

func (p *pathResolver) Add(path string, handler func(res http.ResponseWriter, req *http.Request, server *Server)) {
	p.handlers[path] = handler
}

func (p *pathResolver) Get(pathCheck string) (func(res http.ResponseWriter, req *http.Request, server *Server), error) {
	for pattern, handlerFunc := range p.handlers {
		if ok, err := path.Match(pattern, pathCheck); ok && err == nil {
			return handlerFunc, nil
		} else if err != nil {
			return nil, fmt.Errorf("pathResolver: get: %w", err)
		}
	}

	return nil, nil
}
