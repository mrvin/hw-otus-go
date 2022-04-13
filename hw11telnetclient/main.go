package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mrvin/hw-otus-go/hw11_telnet_client/telnet"
)

func usage() {
	fmt.Printf("usage: %s -host hostname -port port -timeout timeout\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var port int
	var host string
	var timeout time.Duration

	flag.IntVar(&port, "port", 8080, "port")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
	flag.StringVar(&host, "host", "localhost", "dns name or ip address")
	flag.Usage = usage
	flag.Parse()

	confHost := fmt.Sprintf("%s:%d", host, port)
	fmt.Println(confHost)
	client := telnet.NewClient(confHost, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalf("telnetclient: %v", err)
	}

	done := make(chan struct{})
	go func() {
		if err := client.Send(); err != nil {
			log.Print(err)
		}
		done <- struct{}{}
	}()

	if err := client.Receive(); err != nil {
		log.Print(err)
	}
	<-done
}
