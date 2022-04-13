package telnet

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type client struct {
	conn    *net.TCPConn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

type Client interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

func NewClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) Client {
	return &client{address: address,
		timeout: timeout,
		in:      in,
		out:     out}
}

func (c *client) Connect() error {
	var err error
	addr, err := net.ResolveTCPAddr("tcp", c.address)
	if err != nil {
		return fmt.Errorf("can't create: %v", err)
	}

	c.conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return fmt.Errorf("can't connect: %v", err)
	}
	log.Printf("...Connected to %s\n", c.address)

	return nil
}

func (c *client) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("can't close: %v", err)
	}

	return nil
}

// Завершит копирование при нажатии <Ctrl+D>
func (c *client) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("can't send: %v", err)
	}
	log.Print("...EOF")
	c.conn.CloseWrite()

	return nil
}

func (c *client) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("can't receive: %v", err)
	}
	log.Print("...Connection was closed by peer")
	c.conn.CloseRead()

	return nil
}
