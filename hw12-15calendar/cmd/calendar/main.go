//go:generate protoc -I=../../api/ --go_out=../../api/ --go-grpc_out=../../api/ ../../api/eventservice.proto

package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/cmd/calendar/config"
	grpcserver "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/server/grpc"
	httpserver "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/server/http"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	memorystorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/memory"
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
)

var ctx = context.Background()

func main() {
	configFile := flag.String("config", "/etc/calendar/config.yml", "path to configuration file")
	flag.Parse()

	config, err := config.Parse(*configFile)
	if err != nil {
		log.Fatalf("Parse config: %v", err)
	}

	logFile := logInit(&config.Logger)

	var storage storage.Storage
	if config.InMem {
		log.Println("Storage in memory")
		storage = memorystorage.New()
	} else {
		log.Println("Storage in sql")
		storage, err = sqlstorage.New(ctx, &config.DB)
		if err != nil {
			log.Fatalf("db: %v", err)
		}
		log.Println("Connect db")
	}

	serverHTTP := httpserver.New(&config.HTTP, storage)
	serverGRPC, err := grpcserver.New(&config.GRPC, storage)
	if err != nil {
		log.Fatalf("GRPC: %v", err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT /*(Control-C)*/, syscall.SIGTERM)
	go listenForShutdown(signals, serverHTTP, serverGRPC)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := serverHTTP.Start(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("HTTP server: failed to start: %v", err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if err := serverGRPC.Start(); err != nil {
			log.Fatalf("GRPC server: failed to start: %v", err)
		}
	}()

	wg.Wait()

	if storageSQL, ok := storage.(*sqlstorage.Storage); ok {
		log.Println("Close sql storage")
		storageSQL.Close()
	}

	log.Println("Stop service calendar")
	if logFile != nil {
		logFile.Close()
	}
}

func listenForShutdown(signals chan os.Signal, serverHTTP *httpserver.Server, serverGRPC *grpcserver.Server) {
	<-signals
	signal.Stop(signals)

	if err := serverHTTP.Stop(ctx); err != nil {
		log.Fatalf("HTTP server: failed to stop: %v", err)
	}

	serverGRPC.Stop()

}

func logInit(config *config.LoggerConf) *os.File {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if config.FilePath == "" {
		return nil
	}

	logFile, err := os.OpenFile(config.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("log init: %v", err)
		return nil
	}
	log.SetOutput(logFile)

	return logFile
}
