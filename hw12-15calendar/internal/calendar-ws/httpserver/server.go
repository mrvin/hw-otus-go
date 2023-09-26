package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-ws/httpserver/handler"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver"
	regexpresolver "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/resolver/regex"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	serv http.Server
	res  resolver.Resolver
}

type logger struct {
	Inner http.Handler
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	totalWritByte int
}

func New(conf *Conf, grpcclient calendarapi.EventServiceClient) *Server {

	res := regexpresolver.New()
	h := handler.New(grpcclient)
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
	}
	loggerServer := logger{Inner: otelhttp.NewHandler(&server, "HTTP")}

	//nolint:exhaustivestruct,exhaustruct
	server.serv = http.Server{
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler:      &loggerServer,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &server
}

func (s *Server) Start() error {
	slog.Info("Start http server: http://" + s.serv.Addr)
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

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(slByte []byte) (writeByte int, err error) {
	writeByte, err = lrw.ResponseWriter.Write(slByte)
	lrw.totalWritByte += writeByte
	return
}

func (l *logger) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	logReq := slog.With(
		slog.String("method", req.Method),
		slog.String("path", req.URL.Path),
		slog.String("remote_addr", req.RemoteAddr),
		slog.String("user_agent", req.UserAgent()),
		//		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	timeStart := time.Now()
	lrw := NewLoggingResponseWriter(res)
	defer func() {
		logReq.Info("Request "+req.Proto,
			slog.Int("status", lrw.statusCode),
			slog.Int("bytes", lrw.totalWritByte),
			slog.String("duration", time.Since(timeStart).String()),
		)
	}()

	l.Inner.ServeHTTP(lrw, req)
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("Stop http server")
	if err := s.serv.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}
