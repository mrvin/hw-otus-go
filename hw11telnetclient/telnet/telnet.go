package telnet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type ClientInterface interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func NewClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) *Client {
	return &Client{
		conn:    nil,
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *Client) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("can't connect: %w", err)
	}
	log.Printf("...Connected to %s\n", c.address)

	return nil
}

func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("can't close: %w", err)
	}

	return nil
}

// Завершит копирование при нажатии <Ctrl+D>.
func (c *Client) Send() error {
	_, err := io.Copy(c.conn, c.in)
	if err != nil {
		if errors.Is(err, io.ErrClosedPipe) {
			return nil
		}
		return fmt.Errorf("can't send: %w", err)
	}

	log.Print("...EOF")

	return nil
}

func (c *Client) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("can't receive: %w", err)
	}
	log.Print("...Connection was closed by peer")

	return nil
}
