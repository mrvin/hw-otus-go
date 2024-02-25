package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-api"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
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

func New(ctx context.Context, conf *Conf) (*Client, error) {
	var client Client

	ctx, cancel := context.WithTimeout(ctx, shortDuration)
	defer cancel()
	address := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	var err error
	client.conn, err = grpc.DialContext(
		ctx,
		address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", address, err)
	}

	client.eventService = calendarapi.NewEventServiceClient(client.conn)
	client.userService = calendarapi.NewUserServiceClient(client.conn)

	return &client, nil
}

func (c *Client) Close() error {
	c.conn.Close()

	return nil
}
