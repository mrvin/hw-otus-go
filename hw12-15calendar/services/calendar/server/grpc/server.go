package internalgrpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

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
		return nil, fmt.Errorf("establish tcp connection: %v", err)
	}
	server.serv = grpc.NewServer()
	apipb.RegisterEventsServer(server.serv, &server)

	return &server, nil
}

func (s *Server) Start() error {
	log.Print("Start gRPC server")
	if err := s.serv.Serve(s.ln); err != nil {
		return fmt.Errorf("start grpc server: %v", err)
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

func (s *Server) UpdateUser(ctx context.Context, userpb *apipb.User) (*apipb.Null, error) {
	user := storage.User{ID: int(userpb.GetId()), Name: userpb.GetName(), Email: userpb.GetEmail()}
	err := s.stor.UpdateUser(ctx, &user)

	return nil, err
}

func (s *Server) DeleteUser(ctx context.Context, req *apipb.UserRequest) (*apipb.Null, error) {
	err := s.stor.DeleteUser(ctx, int(req.GetId()))

	return nil, err
}

func (s *Server) CreateEvent(ctx context.Context, eventpb *apipb.Event) (*apipb.EventResponse, error) {
	event := storage.Event{Title: eventpb.GetTitle(), Description: eventpb.GetDescription(), /*, StartTime, StopTime*/
		UserID: int(eventpb.GetUserID())}
	err := s.stor.CreateEvent(ctx, &event)

	return &apipb.EventResponse{Id: int64(event.ID)}, err
}

func (s *Server) GetEvent(ctx context.Context, req *apipb.EventRequest) (*apipb.Event, error) {
	event, err := s.stor.GetEvent(ctx, int(req.GetId()))

	return &apipb.Event{Id: int64(event.ID), Title: event.Title, Description: event.Description, UserID: int64(event.UserID)}, err
}

func (s *Server) UpdateEvent(ctx context.Context, eventpb *apipb.Event) (*apipb.Null, error) {
	event := storage.Event{ID: int(eventpb.GetId()), Title: eventpb.GetTitle(),
		Description: eventpb.GetDescription(), /*, StartTime, StopTime*/
		UserID:      int(eventpb.GetUserID())}
	err := s.stor.UpdateEvent(ctx, &event)

	return nil, err
}

func (s *Server) DeleteEvent(ctx context.Context, req *apipb.EventRequest) (*apipb.Null, error) {
	err := s.stor.DeleteEvent(ctx, int(req.GetId()))

	return nil, err
}
