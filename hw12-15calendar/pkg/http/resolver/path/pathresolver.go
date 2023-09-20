package path

import (
	"log/slog"
	"net/http"
	"path"
)

// PathResolver is not thread safe.
type PathResolver struct {
	handlers map[string]http.HandlerFunc
}

func New() *PathResolver {
	return &PathResolver{make(map[string]http.HandlerFunc)}
}

func (p *PathResolver) Add(path string, handler http.HandlerFunc) {
	p.handlers[path] = handler
}

func (p *PathResolver) Get(pathCheck string) http.HandlerFunc {
	for pattern, handlerFunc := range p.handlers {
		if ok, err := path.Match(pattern, pathCheck); ok && err == nil {
			return handlerFunc
		} else if err != nil {
			slog.Error("pathResolver: get: " + err.Error())
		}
	}

	return nil
}
