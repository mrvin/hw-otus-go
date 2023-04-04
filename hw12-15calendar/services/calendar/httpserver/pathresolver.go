package httpserver

import (
	"net/http"
	"path"

	"go.uber.org/zap"
)

// TODO:add thread safety.
type pathResolver struct {
	handlers map[string]http.HandlerFunc
}

func newPathResolver() *pathResolver {
	return &pathResolver{make(map[string]http.HandlerFunc)}
}

func (p *pathResolver) Add(path string, handler http.HandlerFunc) {
	p.handlers[path] = handler
}

func (p *pathResolver) Get(pathCheck string) http.HandlerFunc {
	for pattern, handlerFunc := range p.handlers {
		if ok, err := path.Match(pattern, pathCheck); ok && err == nil {
			return handlerFunc
		} else if err != nil {
			zap.S().Errorf("pathResolver: get: %v", err)
		}
	}

	return nil
}
