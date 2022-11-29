package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/app"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	serv *grpc.Server
	ln   net.Listener
	app  *app.App
	addr string
	log  *zap.SugaredLogger
}

func New(conf *Conf, app *app.App) (*Server, error) {
	var server Server

	server.app = app
	server.log = zap.S()

	var err error
	server.addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	server.ln, err = net.Listen("tcp", server.addr)
	if err != nil {
		return nil, fmt.Errorf("establish tcp connection: %w", err)
	}
	server.serv = grpc.NewServer()
	calendarapi.RegisterEventServiceServer(server.serv, &server)

	return &server, nil
}

func (s *Server) Start() error {
	s.log.Infof("Start gRPC server: %s", s.addr)
	if err := s.serv.Serve(s.ln); err != nil {
		return fmt.Errorf("start grpc server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.log.Info("Stop gRPC server")
	s.serv.GracefulStop()
	s.ln.Close()
}

func (s *Server) CreateUser(ctx context.Context, userpb *calendarapi.User) (*calendarapi.UserResponse, error) {
	defer logGRPC(ctx, "create user")()
	user := storage.User{ID: 0, Name: userpb.GetName(), Email: userpb.GetEmail(), Events: nil}
	if err := s.app.CreateUser(ctx, &user); err != nil {
		err := fmt.Errorf("create user: %w", err)
		s.log.Error(err)
		return nil, err
	}

	return &calendarapi.UserResponse{Id: int64(user.ID)}, nil
}

func (s *Server) GetUser(ctx context.Context, req *calendarapi.UserRequest) (*calendarapi.User, error) {
	defer logGRPC(ctx, "get user")()
	user, err := s.app.GetUser(ctx, int(req.GetId()))
	if err != nil {
		err := fmt.Errorf("get user: %w", err)
		s.log.Error(err)
		return nil, err
	}

	return &calendarapi.User{Id: int64(user.ID), Name: user.Name, Email: user.Email}, nil
}

func (s *Server) GetAllUsers(ctx context.Context, _ *emptypb.Empty) (*calendarapi.Users, error) {
	defer logGRPC(ctx, "get all user")()
	users, err := s.app.GetAllUsers(ctx)
	if err != nil {
		err := fmt.Errorf("get all users: %w", err)
		s.log.Error(err)
		return nil, err
	}

	pbUsers := make([]*calendarapi.User, len(users))
	for i, user := range users {
		pbUsers[i] = &calendarapi.User{Id: int64(user.ID), Name: user.Name, Email: user.Email}
	}

	return &calendarapi.Users{Users: pbUsers}, nil
}

func (s *Server) UpdateUser(ctx context.Context, userpb *calendarapi.User) (*emptypb.Empty, error) {
	defer logGRPC(ctx, "update user")()
	user := storage.User{ID: int(userpb.GetId()), Name: userpb.GetName(), Email: userpb.GetEmail(), Events: nil}
	if err := s.app.UpdateUser(ctx, &user); err != nil {
		err := fmt.Errorf("update user: %w", err)
		s.log.Error(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *calendarapi.UserRequest) (*emptypb.Empty, error) {
	defer logGRPC(ctx, "delete user")()
	if err := s.app.DeleteUser(ctx, int(req.GetId())); err != nil {
		err := fmt.Errorf("delete user: %w", err)
		s.log.Error(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) CreateEvent(ctx context.Context, pbEvent *calendarapi.Event) (*calendarapi.EventResponse, error) {
	defer logGRPC(ctx, "create event")()
	event, err := convertpbEventToEvent(pbEvent)
	if err != nil {
		err := fmt.Errorf("create event: %w", err)
		s.log.Error(err)
		return nil, err
	}
	if err := s.app.CreateEvent(ctx, event); err != nil {
		err := fmt.Errorf("create event: %w", err)
		s.log.Error(err)
		return nil, err
	}

	return &calendarapi.EventResponse{Id: int64(event.ID)}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *calendarapi.EventRequest) (*calendarapi.Event, error) {
	defer logGRPC(ctx, "get event")()
	event, err := s.app.GetEvent(ctx, int(req.GetId()))
	if err != nil {
		err := fmt.Errorf("get event: %w", err)
		s.log.Error(err)
		return nil, err
	}

	return &calendarapi.Event{Id: int64(event.ID), Title: event.Title, Description: event.Description,
		StartTime: timestamppb.New(event.StartTime), StopTime: timestamppb.New(event.StopTime), UserID: int64(event.UserID)}, nil
}

func (s *Server) GetEventsForUser(ctx context.Context, req *calendarapi.GetEventsForUserRequest) (*calendarapi.Events, error) {
	defer logGRPC(ctx, "get event for user")()

	if err := req.DaysAhead.Date.CheckValid(); err != nil {
		return nil, fmt.Errorf("incorrect value date: %w", err)
	}
	date := req.DaysAhead.Date.AsTime()

	events, err := s.app.GetEventsForUser(ctx, int(req.User.GetId()), date, int(req.DaysAhead.Days))
	if err != nil {
		err := fmt.Errorf("get events for user: %w", err)
		s.log.Error(err)
		return nil, err
	}

	pbEvents := make([]*calendarapi.Event, len(events))
	for i, event := range events {
		pbEvents[i] = &calendarapi.Event{Id: int64(event.ID), Title: event.Title, Description: event.Description,
			StartTime: timestamppb.New(event.StartTime), StopTime: timestamppb.New(event.StopTime), UserID: int64(event.UserID)}
	}

	return &calendarapi.Events{Events: pbEvents}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, pbEvent *calendarapi.Event) (*emptypb.Empty, error) {
	defer logGRPC(ctx, "update event")()
	event, err := convertpbEventToEvent(pbEvent)
	if err != nil {
		err := fmt.Errorf("update event: %w", err)
		s.log.Error(err)
		return nil, err
	}

	if err := s.app.UpdateEvent(ctx, event); err != nil {
		err := fmt.Errorf("update event: %w", err)
		s.log.Error(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *calendarapi.EventRequest) (*emptypb.Empty, error) {
	defer logGRPC(ctx, "delete event")()
	if err := s.app.DeleteEvent(ctx, int(req.GetId())); err != nil {
		err := fmt.Errorf("delete event: %w", err)
		s.log.Error(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func convertpbEventToEvent(pbEvent *calendarapi.Event) (*storage.Event, error) {
	if err := pbEvent.StartTime.CheckValid(); err != nil {
		return nil, fmt.Errorf("incorrect value StartTime: %w", err)
	}
	if err := pbEvent.StopTime.CheckValid(); err != nil {
		return nil, fmt.Errorf("incorrect value StopTime: %w", err)
	}
	startTime := pbEvent.StartTime.AsTime()
	stopTime := pbEvent.StopTime.AsTime()

	event := storage.Event{ID: int(pbEvent.GetId()), Title: pbEvent.GetTitle(),
		Description: pbEvent.GetDescription(), StartTime: startTime, StopTime: stopTime,
		UserID: int(pbEvent.GetUserID())}

	return &event, nil
}

func logGRPC(ctx context.Context, operation string) func() {
	start := time.Now()
	log := zap.S()
	ip := ""
	peer, ok := peer.FromContext(ctx)
	if !ok {
		log.Warn("cant get perr")
	} else {
		ip = peer.Addr.String()
	}
	return func() {
		log.Infow("", "ip", ip,
			"operation", operation,
			"proto", "gRPC",
			"duration", time.Since(start))
	}
}
