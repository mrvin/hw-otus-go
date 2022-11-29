package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

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
	log  *zap.SugaredLogger
}

func New(conf *Conf, app *app.App) *Server {
	var server Server

	server.app = app
	server.log = zap.S()
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
	s.log.Infof("Start server: http://%s", s.serv.Addr)
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
		s.log.Error("Get handler function: %v", err)
	} else {
		if handlerFunc != nil {
			handlerFunc(res, req, s)
			return
		}
	}

	http.NotFound(res, req)
}

func logReq(req *http.Request) func() {
	start := time.Now()
	log := zap.S()
	return func() {
		log.Infow("", "ip", strings.Split(req.RemoteAddr, ":")[0],
			"method", req.Method,
			"path", req.URL.Path,
			"proto", req.Proto,
			"duration", time.Since(start) /*, req.Header["User-Agent"]*/)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Stop http server")
	if err := s.serv.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}
