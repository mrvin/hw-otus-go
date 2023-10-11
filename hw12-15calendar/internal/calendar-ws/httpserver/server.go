package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-ws/client"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-ws/httpserver/handler"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver"
	regexpresolver "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver/regex"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	http.Server
}

func New(conf *Conf, client client.Calendar) *Server {
	res := regexpresolver.New()

	h := handler.New(client)

	res.Add("GET /list-users", h.DisplayListUsers)
	res.Add(`GET \/list-events\?id=([0-9]+)\&days=([0-9]+$)`, h.DisplayListEventsForUser)

	res.Add(`GET \/user\?id=([0-9]+$)`, h.DisplayUser)
	res.Add(`GET \/event\?id=([0-9]+$)`, h.DisplayEvent)

	res.Add(`GET \/delete-user\?id=([0-9]+$)`, h.DeleteUser)
	res.Add(`GET \/delete-event\?id=([0-9]+$)`, h.DeleteEvent)

	res.Add("GET /form-user", h.DisplayFormUser)
	res.Add("POST /create-user", h.CreateUser)

	res.Add(`GET \/form-event\?id=([0-9]+$)`, h.DisplayFormEvent)
	res.Add("POST /create-event", h.CreateEvent)

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
	check := req.Method + " " + req.URL.Path + "?" + req.URL.RawQuery
	if handlerFunc := r.Get(check); handlerFunc != nil {
		handlerFunc(res, req)
		return
	}

	http.NotFound(res, req)
}
