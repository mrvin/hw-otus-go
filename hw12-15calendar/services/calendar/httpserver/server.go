package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	serv http.Server
	stor storage.Storage
	pr   *pathResolver
}

func New(conf *Conf, stor storage.Storage) *Server {
	var server Server

	server.stor = stor
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
	log.Print("Start http server")
	if err := s.serv.ListenAndServe(); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	return nil
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer logReq(req)()
	check := req.Method + " " + req.URL.Path
	for pattern, handlerFunc := range s.pr.handlers {
		if ok, err := path.Match(pattern, check); ok && err == nil {
			handlerFunc(res, req, s)

			return
		} else if err != nil {
			fmt.Fprint(res, err)
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