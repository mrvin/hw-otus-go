package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	serv       http.Server
	pr         *regexResolver
	templates  *templateLoader
	grpcclient calendarapi.EventServiceClient
}

func New(conf *Conf, grpcclient calendarapi.EventServiceClient) *Server {
	var server Server

	server.pr = newRegexResolver()
	server.templates = newTemplateLoader()
	server.templates.LoadTemplates("templates")
	server.grpcclient = grpcclient

	server.pr.Add("GET /list-users", displayListUsers)
	server.pr.Add(`GET /list-events\?id=([0-9]+$)`, displayListEventsForUser)

	server.pr.Add(`GET /user\?id=([0-9]+$)`, displayUser)
	server.pr.Add(`GET /event\?id=([0-9]+$)`, displayEvent)

	server.pr.Add(`GET /delete-user\?id=([0-9]+$)`, deleteUser)
	server.pr.Add(`GET /delete-event\?id=([0-9]+$)`, deleteEvent)

	server.pr.Add("GET /form-user", displayFormUser)
	server.pr.Add("POST /create-user", createUser)

	server.pr.Add(`GET /form-event\?id=([0-9]+$)`, displayFormEvent)
	server.pr.Add("POST /create-event", createEvent)

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

	check := req.Method + " " + req.URL.Path + "?" + req.URL.RawQuery
	handlerFunc := s.pr.Get(check)
	if handlerFunc != nil {
		//TODO: Add running in goroutines
		handlerFunc(res, req, s)
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
	log.Print("Stop http server")
	if err := s.serv.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}
