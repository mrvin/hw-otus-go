package httpserver

import (
	"net/http"
	"path"

	"go.uber.org/zap"
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

func (p *pathResolver) Get(pathCheck string) func(res http.ResponseWriter, req *http.Request, server *Server) {
	for pattern, handlerFunc := range p.handlers {
		if ok, err := path.Match(pattern, pathCheck); ok && err == nil {
			return handlerFunc
		} else if err != nil {
			zap.S().Errorf("pathResolver: get: %v", err)
		}
	}

	return nil
}
