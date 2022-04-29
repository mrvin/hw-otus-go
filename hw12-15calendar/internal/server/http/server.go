package internalhttp

import (
	"fmt"
	"net/http"
	"path"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/app"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
)

type Server struct {
	serv http.Server
	logg *logger.Logger
	stor app.Storage
	pr   *pathResolver
}

func New(conf *config.HTTPConf, logg *logger.Logger, stor app.Storage) *Server {
	var server Server

	server.logg = logg
	server.stor = stor
	server.pr = newPathResolver()

	server.pr.Add("GET /hello", hello)

	server.serv = http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler: &server,
	}

	return &server
}

func (s *Server) Start() error {
	if err := s.serv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

func (s *Server) Stop() error {
	if err := s.serv.Shutdown(nil); err != nil {
		return err
	}
	return nil
}
