package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver"
	regexpresolver "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver/regex"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/cmd/calendar-ws/httpserver/handler"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	serv http.Server
	res  resolver.Resolver
	log  *zap.SugaredLogger
}

func New(conf *Conf, grpcclient calendarapi.EventServiceClient) *Server {
	log := zap.S()

	res := regexpresolver.New()
	h := handler.New(grpcclient, log)
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

	server := Server{
		res: res,
		log: log,
	}

	//nolint:exhaustivestruct,exhaustruct
	server.serv = http.Server{
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler:      otelhttp.NewHandler(http.Handler(&server), "HTTP"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
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
	check := req.Method + " " + req.URL.Path + "?" + req.URL.RawQuery
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
