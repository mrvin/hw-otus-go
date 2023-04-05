package resolver

import (
	"net/http"
)

type Resolver interface {
	Add(regex string, handler http.HandlerFunc)
	Get(pathCheck string) http.HandlerFunc
}
