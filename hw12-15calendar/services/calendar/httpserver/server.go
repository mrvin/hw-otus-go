package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/app"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/httpserver/handler"
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
	pr   *pathResolver
	log  *zap.SugaredLogger
}

func New(conf *Conf, app *app.App) *Server {
	var server Server

	log := zap.S()
	server.log = log
	server.pr = newPathResolver()

	h := handler.New(app, log)

	server.pr.Add("POST /users", h.CreateUser)
	server.pr.Add("GET /users", h.GetUser)
	server.pr.Add("PUT /users", h.UpdateUser)
	server.pr.Add("DELETE /users", h.DeleteUser)

	server.pr.Add("POST /events", h.CreateEvent)
	server.pr.Add("GET /events", h.GetEvent)
	server.pr.Add("PUT /events", h.UpdateEvent)
	server.pr.Add("DELETE /events", h.DeleteEvent)

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

func (s *Server) StartTLS(conf *ConfHTTPS) error {
	s.log.Infof("Start server: https://%s", s.serv.Addr)
	if err := s.serv.ListenAndServeTLS(conf.CertFile, conf.KeyFile); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	return nil
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer logReq(req)()

	check := req.Method + " " + req.URL.Path
	if handlerFunc := s.pr.Get(check); handlerFunc != nil {
		handlerFunc(res, req)
		return
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
