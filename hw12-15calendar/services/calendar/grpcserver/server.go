package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	apipb "github.com/mrvin/hw-otus-go/hw12-15calendar/api"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	serv *grpc.Server
	ln   net.Listener
	stor storage.Storage
}

func New(conf *Conf, stor storage.Storage) (*Server, error) {
	var server Server

	server.stor = stor

	var err error
	server.ln, err = net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
	if err != nil {
		return nil, fmt.Errorf("establish tcp connection: %w", err)
	}
	server.serv = grpc.NewServer()
	apipb.RegisterEventsServer(server.serv, &server)

	return &server, nil
}

func (s *Server) Start() error {
	log.Print("Start gRPC server")
	if err := s.serv.Serve(s.ln); err != nil {
		return fmt.Errorf("start grpc server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	log.Print("Stop gRPC server")
	s.serv.GracefulStop()
	s.ln.Close()
}

func (s *Server) CreateUser(ctx context.Context, userpb *apipb.User) (*apipb.UserResponse, error) {
	user := storage.User{Name: userpb.GetName(), Email: userpb.GetEmail()}
	err := s.stor.CreateUser(ctx, &user)

	return &apipb.UserResponse{Id: int64(user.ID)}, err
}

func (s *Server) GetUser(ctx context.Context, req *apipb.UserRequest) (*apipb.User, error) {
	user, err := s.stor.GetUser(ctx, int(req.GetId()))

	return &apipb.User{Id: int64(user.ID), Name: user.Name, Email: user.Email}, err
}

func (s *Server) GetAllUsers(ctx context.Context, null *apipb.Null) (*apipb.Users, error) {
	users, err := s.stor.GetAllUsers(ctx)
	if err != nil {
		err := fmt.Errorf("get all users: %w", err)
		log.Print(err)
		return nil, err
	}

	var pbUsers []*apipb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &apipb.User{Id: int64(user.ID), Name: user.Name, Email: user.Email})
	}

	return &apipb.Users{Users: pbUsers}, err
}

func (s *Server) UpdateUser(ctx context.Context, userpb *apipb.User) (*apipb.Null, error) {
	user := storage.User{ID: int(userpb.GetId()), Name: userpb.GetName(), Email: userpb.GetEmail()}
	err := s.stor.UpdateUser(ctx, &user)

	return nil, err
}

func (s *Server) DeleteUser(ctx context.Context, req *apipb.UserRequest) (*apipb.Null, error) {
	err := s.stor.DeleteUser(ctx, int(req.GetId()))

	return nil, err
}

func (s *Server) CreateEvent(ctx context.Context, pbEvent *apipb.Event) (*apipb.EventResponse, error) {
	event, err := convertpbEventToEvent(pbEvent)
	if err != nil {
		err := fmt.Errorf("create event: %w", err)
		log.Print(err)
		return nil, err
	}
	if err := s.stor.CreateEvent(ctx, event); err != nil {
		err := fmt.Errorf("create event: %w", err)
		log.Print(err)
		return nil, err
	}

	return &apipb.EventResponse{Id: int64(event.ID)}, err
}

func (s *Server) GetEvent(ctx context.Context, req *apipb.EventRequest) (*apipb.Event, error) {
	event, err := s.stor.GetEvent(ctx, int(req.GetId()))

	return &apipb.Event{Id: int64(event.ID), Title: event.Title, Description: event.Description,
		StartTime: timestamppb.New(event.StartTime), StopTime: timestamppb.New(event.StopTime), UserID: int64(event.UserID)}, err
}

func (s *Server) GetEventsForUser(ctx context.Context, req *apipb.UserRequest) (*apipb.Events, error) {
	events, err := s.stor.GetEventsForUser(ctx, int(req.GetId()))
	if err != nil {
		err := fmt.Errorf("get events for user: %w", err)
		log.Print(err)
		return nil, err
	}

	var pbEvents []*apipb.Event
	for _, event := range events {
		pbEvents = append(pbEvents, &apipb.Event{Id: int64(event.ID), Title: event.Title, Description: event.Description,
			StartTime: timestamppb.New(event.StartTime), StopTime: timestamppb.New(event.StopTime)})
	}

	return &apipb.Events{Events: pbEvents}, err
}

func (s *Server) UpdateEvent(ctx context.Context, pbEvent *apipb.Event) (*apipb.Null, error) {
	event, err := convertpbEventToEvent(pbEvent)
	if err != nil {
		err := fmt.Errorf("update event: %w", err)
		log.Print(err)
		return nil, err
	}

	if err := s.stor.UpdateEvent(ctx, event); err != nil {
		err := fmt.Errorf("update event: %w", err)
		log.Print(err)
		return nil, err
	}

	return &apipb.Null{}, err
}

func (s *Server) DeleteEvent(ctx context.Context, req *apipb.EventRequest) (*apipb.Null, error) {
	err := s.stor.DeleteEvent(ctx, int(req.GetId()))

	return nil, err
}

func convertpbEventToEvent(pbEvent *apipb.Event) (*storage.Event, error) {
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
