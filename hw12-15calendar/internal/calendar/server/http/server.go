package httpserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	handlerevent "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers/event"
	handleruserdelete "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers/user/delete"
	handleruserget "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers/user/get"
	handlerusersignin "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers/user/signin"
	handlerusersignup "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers/user/signup"
	handleruserupdate "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers/user/update"
	authservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/auth"
	eventservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/event"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const readTimeout = 5   // in second
const writeTimeout = 10 // in second
const idleTimeout = 1   // in minute

//nolint:tagliatelle
type ConfHTTPS struct {
	CertFile       string `yaml:"cert_file"`
	KeyFile        string `yaml:"key_file"`
	ClientCertFile string `yaml:"client_cert_file"`
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
	mux := http.NewServeMux()

	handlerEvent := handlerevent.New(events)

	mux.HandleFunc(http.MethodPost+" /signup", handlerusersignup.New(auth))
	mux.HandleFunc(http.MethodGet+" /login", handlerusersignin.New(auth))

	mux.HandleFunc(http.MethodGet+" /user", auth.Authorized(handleruserget.New(auth)))
	mux.HandleFunc(http.MethodPut+" /user", auth.Authorized(handleruserupdate.New(auth)))
	mux.HandleFunc(http.MethodDelete+" /user", auth.Authorized(handleruserdelete.New(auth)))

	mux.HandleFunc(http.MethodPost+" /event", auth.Authorized(handlerEvent.CreateEvent))
	mux.HandleFunc(http.MethodGet+" /event", auth.Authorized(handlerEvent.GetEvent))
	mux.HandleFunc(http.MethodPut+" /event", auth.Authorized(handlerEvent.UpdateEvent))
	mux.HandleFunc(http.MethodDelete+" /event", auth.Authorized(handlerEvent.DeleteEvent))

	loggerServer := logger.Logger{Inner: otelhttp.NewHandler(mux, "HTTP")}

	var tlsConf *tls.Config
	if conf.IsHTTPS {
		clientCert, err := os.ReadFile(conf.HTTPS.ClientCertFile)
		if err != nil {
			slog.Warn("Failed to read client tls certificate file: " + err.Error())
		} else {
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(clientCert)

			//nolint:exhaustruct
			tlsConf = &tls.Config{
				ClientCAs:  pool,
				ClientAuth: tls.RequireAndVerifyClientCert,
				MinVersion: tls.VersionTLS12,
			}
		}
	}

	return &Server{
		//nolint:exhaustivestruct,exhaustruct
		http.Server{
			Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			Handler:      &loggerServer,
			ReadTimeout:  readTimeout * time.Second,
			WriteTimeout: writeTimeout * time.Second,
			IdleTimeout:  idleTimeout * time.Minute,
			TLSConfig:    tlsConf,
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
