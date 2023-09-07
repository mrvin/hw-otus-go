package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver"
	pathresolver "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver/path"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/app"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/cmd/calendar/httpserver/handler"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

//nolint:tagliatelle
type ConfHTTPS struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

//nolint:tagliatelle
type Conf struct {
	Host    string    `yaml:"host"`
	Port    int       `yaml:"port"`
	IsHTTPS bool      `yaml:"is_https"`
	HTTPS   ConfHTTPS `yaml:"https"`
}

type Server struct {
	serv http.Server
	res  resolver.Resolver
	log  *zap.SugaredLogger
}

func New(conf *Conf, app *app.App) *Server {
	var server Server

	log := zap.S()
	server.log = log
	server.res = pathresolver.New()

	h := handler.New(app, log)

	server.res.Add("POST /users", h.CreateUser)
	server.res.Add("GET /users", h.GetUser)
	server.res.Add("PUT /users", h.UpdateUser)
	server.res.Add("DELETE /users", h.DeleteUser)

	server.res.Add("POST /events", h.CreateEvent)
	server.res.Add("GET /events", h.GetEvent)
	server.res.Add("PUT /events", h.UpdateEvent)
	server.res.Add("DELETE /events", h.DeleteEvent)

	//nolint:exhaustivestruct,exhaustruct
	server.serv = http.Server{
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler:      otelhttp.NewHandler(http.Handler(&server), "HTTP"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
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

func (s *Server) StartTLS(conf *ConfHTTPS) error {
	s.log.Infof("Start server: https://%s", s.serv.Addr)
	if err := s.serv.ListenAndServeTLS(conf.CertFile, conf.KeyFile); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	return nil
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	if handlerFunc := s.res.Get(check); handlerFunc != nil {
		handlerFunc(res, req)
		return
	}

	http.NotFound(res, req)
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Stop http server")
	if err := s.serv.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}
