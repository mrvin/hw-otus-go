package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver"
	regexpresolver "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver/regex"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar-ws/httpserver/handler"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar-ws/httpserver/templateloader"
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
	var server Server

	log := zap.S()
	server.log = log
	server.res = regexpresolver.New()

	h := handler.New(templateloader.New(), grpcclient, log)

	server.res.Add("GET /list-users", h.DisplayListUsers)
	server.res.Add(`GET /list-events\?id=([0-9]+$)`, h.DisplayListEventsForUser)

	server.res.Add(`GET /user\?id=([0-9]+$)`, h.DisplayUser)
	server.res.Add(`GET /event\?id=([0-9]+$)`, h.DisplayEvent)

	server.res.Add(`GET /delete-user\?id=([0-9]+$)`, h.DeleteUser)
	server.res.Add(`GET /delete-event\?id=([0-9]+$)`, h.DeleteEvent)

	server.res.Add("GET /form-user", h.DisplayFormUser)
	server.res.Add("POST /create-user", h.CreateUser)

	server.res.Add(`GET /form-event\?id=([0-9]+$)`, h.DisplayFormEvent)
	server.res.Add("POST /create-event", h.CreateEvent)

	server.serv = http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler: &server,
	}

	return &server
}

func (s *Server) Start() error {
	s.log.Infof("Start http server: %s", s.serv.Addr)
	if err := s.serv.ListenAndServe(); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	return nil
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer logReq(req)()

	check := req.Method + " " + req.URL.Path + "?" + req.URL.RawQuery
	if handlerFunc := s.res.Get(check); handlerFunc != nil {
		handlerFunc(res, req)
		return
	}

	http.NotFound(res, req)
}

func logReq(req *http.Request) func() {
	start := time.Now()

	return func() {
		log.Printf("%s [%s] %s %s %s %s %s", strings.Split(req.RemoteAddr, ":")[0], start.Format(time.ANSIC),
			req.Method, req.URL.Path, req.URL.RawQuery, req.Proto, time.Since(start) /*, req.Header["User-Agent"]*/)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Infof("Stop http server")
	if err := s.serv.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}
