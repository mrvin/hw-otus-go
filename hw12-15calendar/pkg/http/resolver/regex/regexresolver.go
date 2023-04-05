package regexp

import (
	"net/http"
	"regexp"
)

// RegexpResolver is not thread safe.
type RegexpResolver struct {
	handlers map[string]http.HandlerFunc
	cache    map[string]*regexp.Regexp
}

func New() *RegexpResolver {
	return &RegexpResolver{
		handlers: make(map[string]http.HandlerFunc),
		cache:    make(map[string]*regexp.Regexp),
	}
}

func (r *RegexpResolver) Add(regex string, handler http.HandlerFunc) {
	r.handlers[regex] = handler
	r.cache[regex] = regexp.MustCompile(regex)
}

func (r *RegexpResolver) Get(pathCheck string) http.HandlerFunc {
	for pattern, handlerFunc := range r.handlers {
		if r.cache[pattern].MatchString(pathCheck) {
			return handlerFunc
		}
	}

	return nil
}
