package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/app"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/httpserver/handlers"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver"
	pathresolver "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver/path"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	http.Server
}

func New(conf *Conf, app *app.App) *Server {
	res := pathresolver.New()

	h := handler.New(app)

	res.Add("POST /users", h.CreateUser)
	res.Add("GET /users", h.GetUser)
	res.Add("PUT /users", h.UpdateUser)
	res.Add("DELETE /users", h.DeleteUser)

	res.Add("POST /events", h.CreateEvent)
	res.Add("GET /events", h.GetEvent)
	res.Add("PUT /events", h.UpdateEvent)
	res.Add("DELETE /events", h.DeleteEvent)

	loggerServer := logger.Logger{Inner: otelhttp.NewHandler(&Router{res}, "HTTP")}

	return &Server{
		//nolint:exhaustivestruct,exhaustruct
		http.Server{
			Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			Handler:      &loggerServer,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  1 * time.Minute,
		},
	}
}

func (s *Server) Start() error {
	slog.Info("Start http server: http://" + s.Addr)
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	return nil
}

func (s *Server) StartTLS(conf *ConfHTTPS) error {
	slog.Info("Start http server: https://" + s.Addr)
	if err := s.ListenAndServeTLS(conf.CertFile, conf.KeyFile); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("Stop http server")
	if err := s.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}

type Router struct {
	resolver.Resolver
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	if handlerFunc := r.Get(check); handlerFunc != nil {
		handlerFunc(res, req)
		return
	}

	http.NotFound(res, req)
}