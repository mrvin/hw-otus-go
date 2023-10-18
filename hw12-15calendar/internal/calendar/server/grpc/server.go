package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-api"
	authservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/auth"
	eventservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/event"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		),
	)
	calendarapi.RegisterEventServiceServer(server.serv, &server)
	calendarapi.RegisterUserServiceServer(server.serv, &server)

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

func (s *Server) CreateUser(ctx context.Context, userpb *calendarapi.CreateUserRequest) (*calendarapi.CreateUserResponse, error) {
	user := storage.User{
		ID:           0,
		Name:         userpb.GetName(),
		HashPassword: userpb.GetPassword(),
		Email:        userpb.GetEmail(),
		Events:       nil,
	}
	id, err := s.authService.CreateUser(ctx, &user)
	if err != nil {
		err = fmt.Errorf("create user: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &calendarapi.CreateUserResponse{Id: id}, nil
}

func (s *Server) GetUserByID(ctx context.Context, req *calendarapi.GetUserByIDRequest) (*calendarapi.UserResponse, error) {
	user, err := s.authService.GetUser(ctx, req.GetId())
	if err != nil {
		err := fmt.Errorf("get user: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &calendarapi.UserResponse{
		Id:           user.ID,
		Name:         user.Name,
		HashPassword: user.HashPassword,
		Email:        user.Email,
	}, nil
}

func (s *Server) ListUsers(ctx context.Context, _ *emptypb.Empty) (*calendarapi.ListUsersResponse, error) {
	users, err := s.authService.ListUsers(ctx)
	if err != nil {
		err := fmt.Errorf("get all users: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	pbUsers := make([]*calendarapi.UserResponse, len(users))
	for i, user := range users {
		pbUsers[i] = &calendarapi.UserResponse{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return &calendarapi.ListUsersResponse{Users: pbUsers}, nil
}

func (s *Server) UpdateUser(ctx context.Context, userpb *calendarapi.UpdateUserRequest) (*emptypb.Empty, error) {
	user := storage.User{
		Name:   userpb.GetName(),
		Email:  userpb.GetEmail(),
		Events: nil,
	}
	if err := s.authService.UpdateUser(ctx, &user); err != nil {
		err := fmt.Errorf("update user: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *calendarapi.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := s.authService.DeleteUserByName(ctx, req.GetName()); err != nil {
		err := fmt.Errorf("delete user: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) Login(ctx context.Context, req *calendarapi.LoginRequest) (*calendarapi.LoginResponse, error) {
	tokenString, err := s.authService.Authenticate(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &calendarapi.LoginResponse{AccessToken: tokenString}, nil
}

func (s *Server) CreateEvent(ctx context.Context, pbEvent *calendarapi.CreateEventRequest) (*calendarapi.CreateEventResponse, error) {
	if err := pbEvent.StartTime.CheckValid(); err != nil {
		err = fmt.Errorf("incorrect value StartTime: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	if err := pbEvent.StopTime.CheckValid(); err != nil {
		err = fmt.Errorf("incorrect value StopTime: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	event := storage.Event{
		ID:          0,
		Title:       pbEvent.GetTitle(),
		Description: pbEvent.GetDescription(),
		StartTime:   pbEvent.StartTime.AsTime(),
		StopTime:    pbEvent.StopTime.AsTime(),
		UserID:      pbEvent.GetUserID(),
	}

	id, err := s.eventService.CreateEvent(ctx, &event)
	if err != nil {
		err = fmt.Errorf("create event: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &calendarapi.CreateEventResponse{Id: id}, nil
}

func (s *Server) GetEventByID(ctx context.Context, req *calendarapi.GetEventByIDRequest) (*calendarapi.EventResponse, error) {
	event, err := s.eventService.GetEvent(ctx, req.GetId())
	if err != nil {
		err := fmt.Errorf("get event: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &calendarapi.EventResponse{
		Id:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   timestamppb.New(event.StartTime),
		StopTime:    timestamppb.New(event.StopTime),
		UserID:      event.UserID,
	}, nil
}

func (s *Server) ListEventsForUser(ctx context.Context, req *calendarapi.ListEventsForUserRequest) (*calendarapi.ListEventsResponse, error) {
	if err := req.Date.CheckValid(); err != nil {
		return nil, fmt.Errorf("incorrect value date: %w", err)
	}
	date := req.Date.AsTime()

	events, err := s.eventService.ListEventsForUser(ctx, req.GetUserId(), date, int(req.Days))
	if err != nil {
		err := fmt.Errorf("get events for user: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	pbEvents := make([]*calendarapi.EventResponse, len(events), len(events))
	for i, event := range events {
		pbEvents[i] = &calendarapi.EventResponse{
			Id:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   timestamppb.New(event.StartTime),
			StopTime:    timestamppb.New(event.StopTime),
			UserID:      event.UserID,
		}
	}

	return &calendarapi.ListEventsResponse{Events: pbEvents}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, pbEvent *calendarapi.UpdateEventRequest) (*emptypb.Empty, error) {
	if err := pbEvent.StartTime.CheckValid(); err != nil {
		err = fmt.Errorf("incorrect value StartTime: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	if err := pbEvent.StopTime.CheckValid(); err != nil {
		err = fmt.Errorf("incorrect value StopTime: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	event := storage.Event{
		ID:          0,
		Title:       pbEvent.GetTitle(),
		Description: pbEvent.GetDescription(),
		StartTime:   pbEvent.StartTime.AsTime(),
		StopTime:    pbEvent.StopTime.AsTime(),
		UserID:      pbEvent.GetUserID(),
	}

	if err := s.eventService.UpdateEvent(ctx, &event); err != nil {
		err := fmt.Errorf("update event: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *calendarapi.DeleteEventRequest) (*emptypb.Empty, error) {
	if err := s.eventService.DeleteEvent(ctx, req.GetId()); err != nil {
		err := fmt.Errorf("delete event: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
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

/*
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        log.Println("--> unary interceptor: ", info.FullMethod)

        err := interceptor.authorize(ctx, info.FullMethod)
        if err != nil {
            return nil, err
        }

        return handler(ctx, req)
    }
}

type AuthInterceptor struct {
    authClient  *AuthClient
    authMethods map[string]bool
    accessToken string
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
    accessibleRoles, ok := interceptor.accessibleRoles[method]
    if !ok {
        // everyone can access
        return nil
    }

    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return status.Errorf(codes.Unauthenticated, "metadata is not provided")
    }

    values := md["authorization"]
    if len(values) == 0 {
        return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
    }

    accessToken := values[0]
    claims, err := interceptor.jwtManager.Verify(accessToken)
    if err != nil {
        return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
    }

    for _, role := range accessibleRoles {
        if role == claims.Role {
            return nil
        }
    }

    return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}

func (interceptor *AuthInterceptor) refreshToken() error {
    accessToken, err := interceptor.authClient.Login()
    if err != nil {
        return err
    }

    interceptor.accessToken = accessToken
    log.Printf("token refreshed: %v", accessToken)

    return nil
}
*/
