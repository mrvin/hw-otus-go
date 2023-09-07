package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Client struct {
	Cl   calendarapi.EventServiceClient
	conn *grpc.ClientConn
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

	client.Cl = calendarapi.NewEventServiceClient(client.conn)

	return &client, nil
}

func (c *Client) Close() error {
	c.conn.Close()

	return nil
}
