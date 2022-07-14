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

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	memorystorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/memory"
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/grpcserver"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/httpserver"
)

var ctx = context.Background()

func main() {
	configFile := flag.String("config", "/etc/calendar/config.yml", "path to configuration file")
	flag.Parse()

	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		log.Printf("Parse config: %v", err)
		return
	}

	logFile := logger.LogInit(&conf.Logger)

	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	var storage storage.Storage
	if conf.InMem {
		log.Println("Storage in memory")
		storage = memorystorage.New()
	} else {
		var err error
		log.Println("Storage in sql")
		storage, err = sqlstorage.New(ctx, &conf.DB)
		if err != nil {
			log.Printf("db: %v", err)
			return
		}
		log.Println("Connect db")
	}

	serverHTTP := httpserver.New(&conf.HTTP, storage)
	serverGRPC, err := grpcserver.New(&conf.GRPC, storage)
	if err != nil {
		log.Printf("GRPC: %v", err)
		return
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
				log.Printf("HTTP server: failed to start: %v", err)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		if err := serverGRPC.Start(); err != nil {
			log.Printf("GRPC server: failed to start: %v", err)
			return
		}
	}()

	wg.Wait()

	if storageSQL, ok := storage.(*sqlstorage.Storage); ok {
		log.Println("Close sql storage")
		storageSQL.Close()
	}

	log.Println("Stop service calendar")
}

func listenForShutdown(signals chan os.Signal, serverHTTP *httpserver.Server, serverGRPC *grpcserver.Server) {
	<-signals
	signal.Stop(signals)

	if err := serverHTTP.Stop(ctx); err != nil {
		log.Printf("HTTP server: failed to stop: %v", err)
		return
	}

	serverGRPC.Stop()
}
