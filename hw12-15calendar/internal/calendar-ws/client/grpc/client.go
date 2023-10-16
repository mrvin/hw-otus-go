package grpcclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-api"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Client struct {
	eventService calendarapi.EventServiceClient
	userService  calendarapi.UserServiceClient
	conn         *grpc.ClientConn
}

const shortDuration = 5 * time.Second

func New(conf *Conf) (*Client, error) {
	var client Client

	ctx, _ := context.WithTimeout(context.Background(), shortDuration)
	confHost := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	var err error
	client.conn, err = grpc.DialContext(ctx, confHost,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", confHost, err)
	}

	client.eventService = calendarapi.NewEventServiceClient(client.conn)
	client.userService = calendarapi.NewUserServiceClient(client.conn)

	return &client, nil
}
func (c *Client) CreateUser(ctx context.Context, name, password, email string) (int64, error) {
	user := &calendarapi.CreateUserRequest{
		Name:     name,
		Password: password,
		Email:    email,
	}
	response, err := c.userService.CreateUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("gRPC: %w", err)
	}
	slog.Debug("Added user",
		slog.Int64("id", response.Id),
		slog.String("username", user.Name),
		slog.String("password", user.Password),
		slog.String("email", user.Email),
	)

	return response.Id, nil
}

func (c *Client) GetUser(ctx context.Context, id int64) (*storage.User, error) {
	reqUser := &calendarapi.GetUserByIDRequest{Id: id}
	user, err := c.userService.GetUserByID(ctx, reqUser)
	if err != nil {
		return nil, fmt.Errorf("gRPC: %w", err)
	}

	return &storage.User{
		ID:           user.Id,
		Name:         user.Name,
		HashPassword: user.HashPassword,
		Email:        user.Email,
	}, nil
}

func (c *Client) UpdateUser(ctx context.Context, name, password, email string) error {
	user := &calendarapi.UpdateUserRequest{
		Name:     name,
		Password: password,
		Email:    email,
	}
	_, err := c.userService.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("gRPC: %w", err)
	}
	slog.Debug("Update user",
		slog.String("username", name),
	)

	return nil
}

func (c *Client) DeleteUser(ctx context.Context, name string) error {
	reqUser := &calendarapi.DeleteUserRequest{Name: name}
	if _, err := c.userService.DeleteUser(ctx, reqUser); err != nil {
		return fmt.Errorf("gRPC: %w", err)
	}

	return nil
}

func (c *Client) ListUsers(ctx context.Context) ([]storage.User, error) {
	users, err := c.userService.ListUsers(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("gRPC: %w", err)
	}

	listUsers := make([]storage.User, 0, len(users.Users))
	for _, user := range users.Users {
		listUsers = append(listUsers, storage.User{
			ID:           user.Id,
			Name:         user.Name,
			HashPassword: user.HashPassword,
			Email:        user.Email,
		})
	}

	return listUsers, nil
}

func (c *Client) CreateEvent(
	ctx context.Context, title, description string,
	startTime, stopTime time.Time,
	userID int64) (int64, error) {

	event := &calendarapi.CreateEventRequest{
		Title:       title,
		Description: description,
		StartTime:   timestamppb.New(startTime),
		StopTime:    timestamppb.New(stopTime),
		UserID:      userID,
	}
	response, err := c.eventService.CreateEvent(ctx, event)
	if err != nil {
		return 0, fmt.Errorf("gRPC: %w", err)
	}
	slog.Debug("Added event",
		slog.Int64("id", response.Id),
		slog.String("title", event.Title),
		slog.String("Description", event.Description),
	)

	return response.Id, nil
}

func (c *Client) GetEvent(ctx context.Context, id int64) (*storage.Event, error) {
	reqEvent := &calendarapi.GetEventByIDRequest{Id: id}
	event, err := c.eventService.GetEventByID(ctx, reqEvent)
	if err != nil {
		return nil, fmt.Errorf("gRPC: %w", err)
	}

	return &storage.Event{
		ID:          event.Id,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime.AsTime(),
		StopTime:    event.StopTime.AsTime(),
	}, nil
}
func (c *Client) DeleteEvent(ctx context.Context, id int64) error {
	reqEvent := &calendarapi.DeleteEventRequest{Id: id}
	if _, err := c.eventService.DeleteEvent(ctx, reqEvent); err != nil {
		return fmt.Errorf("gRPC: %w", err)

	}

	return nil
}

func (c *Client) ListEventsForUser(ctx context.Context, idUser int64, days int) ([]storage.Event, error) {
	reqUser := &calendarapi.ListEventsForUserRequest{
		UserId: idUser,
		Days:   int32(days),
		Date:   timestamppb.New(time.Now()),
	}

	events, err := c.eventService.ListEventsForUser(ctx, reqUser)
	if err != nil {
		return nil, fmt.Errorf("gRPC: %w", err)
	}
	listEvents := make([]storage.Event, 0, len(events.Events))
	for _, event := range events.Events {
		listEvents = append(listEvents, storage.Event{
			ID:          event.Id,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime.AsTime(),
			StopTime:    event.StopTime.AsTime(),
			UserID:      event.UserID,
		})
	}

	return listEvents, nil
}

func (c *Client) Close() error {
	c.conn.Close()

	return nil
}
