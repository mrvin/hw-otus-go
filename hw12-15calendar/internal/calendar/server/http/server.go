package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	handlerevent "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers/event"
	handleruser "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers/user"
	authservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/auth"
	eventservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/event"
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

func New(conf *Conf, auth *authservice.AuthService, events *eventservice.EventService) *Server {
	res := pathresolver.New()

	handlerEvent := handlerevent.New(events)
	handlerUser := handleruser.New(auth)

	res.Add("POST /signup", handlerUser.SignUp)
	res.Add("GET /login", handlerUser.SignIn)

	res.Add("GET /user", auth.Authorized(handlerUser.GetUser))
	res.Add("PUT /user", auth.Authorized(handlerUser.UpdateUser))
	res.Add("DELETE /user", auth.Authorized(handlerUser.DeleteUser))

	res.Add("POST /event", auth.Authorized(handlerEvent.CreateEvent))
	res.Add("GET /event", auth.Authorized(handlerEvent.GetEvent))
	res.Add("PUT /event", auth.Authorized(handlerEvent.UpdateEvent))
	res.Add("DELETE /event", auth.Authorized(handlerEvent.DeleteEvent))

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
