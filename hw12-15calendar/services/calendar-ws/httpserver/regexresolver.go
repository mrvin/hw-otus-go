package httpserver

import (
	"net/http"
	"regexp"
)

type regexResolver struct {
	handlers map[string]func(res http.ResponseWriter, req *http.Request, server *Server)
	cache    map[string]*regexp.Regexp
}

func newRegexResolver() *regexResolver {
	return &regexResolver{
		handlers: make(map[string]func(res http.ResponseWriter, req *http.Request, server *Server)),
		cache:    make(map[string]*regexp.Regexp),
	}
}

func (r *regexResolver) Add(regex string, handler func(res http.ResponseWriter, req *http.Request, server *Server)) {
	r.handlers[regex] = handler
	r.cache[regex] = regexp.MustCompile(regex)
}

func (r *regexResolver) Get(pathCheck string) func(res http.ResponseWriter, req *http.Request, server *Server) {
	for pattern, handlerFunc := range r.handlers {
		if r.cache[pattern].MatchString(pathCheck) {
			return handlerFunc
		}
	}

	return nil
}
