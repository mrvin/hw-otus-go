package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/app"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	serv http.Server
	app  *app.App
	pr   *pathResolver
}

func New(conf *Conf, app *app.App) *Server {
	var server Server

	server.app = app
	server.pr = newPathResolver()

	server.pr.Add("POST /users", handleCreateUser)
	server.pr.Add("GET /users", handleGetUser)
	server.pr.Add("PUT /users", handleUpdateUser)
	server.pr.Add("DELETE /users", handleDeleteUser)

	server.pr.Add("POST /events", handleCreateEvent)
	server.pr.Add("GET /events", handleGetEvent)
	server.pr.Add("PUT /events", handleUpdateEvent)
	server.pr.Add("DELETE /events", handleDeleteEvent)

	server.serv = http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler: &server,
	}

	return &server
}

func (s *Server) Start() error {
	log.Printf("Start http server: %s", s.serv.Addr)
	if err := s.serv.ListenAndServe(); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	return nil
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer logReq(req)()

	check := req.Method + " " + req.URL.Path
	handlerFunc, err := s.pr.Get(check)
	if err != nil {
		log.Printf("%v", err)
	} else {
		if handlerFunc != nil {
			//TODO: Add running in goroutines
			handlerFunc(res, req, s)
			return
		}
	}

	http.NotFound(res, req)
}

func logReq(req *http.Request) func() {
	start := time.Now()

	return func() {
		log.Printf("%s [%s] %s %s %s %s", strings.Split(req.RemoteAddr, ":")[0], start.Format(time.ANSIC),
			req.Method, req.URL.Path, req.Proto, time.Since(start) /*, req.Header["User-Agent"]*/)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	log.Print("Stop http server")
	if err := s.serv.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}
