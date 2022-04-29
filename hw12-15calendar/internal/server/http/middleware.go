package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func hello(res http.ResponseWriter, req *http.Request, server *Server) {
	data := time.Now()

	fmt.Fprint(res, "Hello world!")

	nanosec := time.Since(data).Nanoseconds()
	server.logg.Printf("%s [%s] %s %s %s %d", req.RemoteAddr, data.Format(time.ANSIC), req.Method, req.URL.Path, req.Proto, nanosec /*, req.Header["User-Agent"]*/)
}
