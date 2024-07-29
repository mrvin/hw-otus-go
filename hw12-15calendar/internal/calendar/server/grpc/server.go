package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	authservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/auth"
	eventservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/event"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/grpcapi"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"google.golang.org/grpc/peer"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	serv         *grpc.Server
	ln           net.Listener
	authService  *authservice.AuthService
	eventService *eventservice.EventService
	addr         string
}

func New(conf *Conf, auth *authservice.AuthService, events *eventservice.EventService) (*Server, error) {
	var server Server

	server.authService = auth
	server.eventService = events

	var err error
	server.addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	server.ln, err = net.Listen("tcp", server.addr)
	if err != nil {
		return nil, fmt.Errorf("establish tcp connection: %w", err)
	}

	server.serv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			LogRequestGRPC,
			server.Auth,
		),
	)
	grpcapi.RegisterEventServiceServer(server.serv, &server)
	grpcapi.RegisterUserServiceServer(server.serv, &server)

	return &server, nil
}

func (s *Server) Start() error {
	slog.Info("Start gRPC server: " + s.addr)
	if err := s.serv.Serve(s.ln); err != nil {
		return fmt.Errorf("start grpc server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	slog.Info("Stop gRPC server")
	s.serv.GracefulStop()
	s.ln.Close()
}

// LogRequest is a gRPC UnaryServerInterceptor that will log the API call to stdOut
func LogRequestGRPC(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (response interface{}, err error) {
	var addr string
	p, ok := peer.FromContext(ctx)
	if !ok {
		slog.Warn("Cant get perr")
	} else {
		addr = p.Addr.String()
	}
	slog.Info("Request gRPC",
		slog.String("addr", addr),
		slog.String("Method", info.FullMethod),
	)
	// Last but super important, execute the handler so that the actualy gRPC request is also performed
	return handler(ctx, req)
}

type TokenGetter interface {
	GetAccessToken() string
}

func (s *Server) Auth(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (response interface{}, err error) {
	if info.FullMethod != "/calendar.UserService/CreateUser" && info.FullMethod != "/calendar.UserService/Login" {
		reqTokenGetter, ok := req.(TokenGetter)
		if !ok {
			panic("cant make request TokenGetter interface")
		}
		tokenString := reqTokenGetter.GetAccessToken()

		claims, err := s.authService.ParseToken(tokenString)
		if err != nil {
			return nil, err
		}
		if info.FullMethod == "/calendar.UserService/ListUsers" {
			role := claims["role"]
			if role != "admin" {
				return nil, fmt.Errorf("not authorization")
			}
		}
		username := claims["username"]
		ctx = context.WithValue(ctx, "username", username)
	}
	return handler(ctx, req)
}

func GetUserName(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if userName, ok := ctx.Value("username").(string); ok {
		return userName
	}
	return ""
}
