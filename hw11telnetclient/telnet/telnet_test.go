package telnet

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		if err != nil {
			t.Fatalf("can't listen: %v", err)
		}
		defer l.Close()

		done := make(chan struct{})
		sendMsg := "hello\n"
		recvMsg := "world\n"
		go func() {
			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			if err != nil {
				t.Fatalf("can't parse timeout: %v", err)
			}

			client := NewClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			if err := client.Connect(); err != nil {
				t.Fatalf("can't connect to host: %v", err)
			}
			defer client.Close()

			in.WriteString(sendMsg)
			if err = client.Send(); err != nil {
				t.Errorf("can't send: %v", err)
			}

			if err := client.Receive(); err != nil {
				t.Errorf("can't receive: %v", err)
			}
			if recvMsg != out.String() {
				t.Errorf("received message was not recorded to out")
			}

			done <- struct{}{}
		}()

		conn, err := l.Accept()
		if err != nil {
			t.Fatalf("can't accept connect: %v", err)
		}

		request := make([]byte, 1024)
		n, err := conn.Read(request)
		if err != nil {
			t.Errorf("can't read sent message: %v", err)
		}
		readMsg := string(request)[:n]
		if sendMsg != readMsg {
			t.Errorf("send message \"%s\" not equal read message \"%s\"", sendMsg, readMsg)
		}

		_, err = conn.Write([]byte(recvMsg))
		if err != nil {
			t.Errorf("can't write received message: %v", err)
		}

		conn.Close()

		<-done
	})
}
