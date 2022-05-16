package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/app"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	httpserver "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/server/http"
	memorystorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/memory"
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
)

var ctx = context.Background()

func main() {
	configFile := flag.String("config", "/etc/calendar/config.yml", "path to configuration file")
	flag.Parse()

	config, err := config.Parse(*configFile)
	if err != nil {
		log.Fatalf("can't parse config: %v", err)
	}

	logg, err := logger.Create(&config.Logger)
	if err != nil {
		log.Fatalf("can't create logger: %v", err)
	}

	logg.Println("Start service calendar")

	var storage app.Storage
	if config.InMem {
		storage = memorystorage.New()
	} else {
		var storageSQL sqlstorage.Storage
		if err := storageSQL.Connect(ctx, &config.DB); err != nil {
			logg.Fatalf("can't connection db: %v", err)
		}
		if err := storageSQL.PrepareQuery(ctx); err != nil {
			logg.Fatalf("can't prepare query: %v", err)
		}

		storage = &storageSQL
	}

	server := httpserver.New(&config.HTTP, logg, storage)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)
	done := make(chan struct{})
	go listenForShutdown(signals, logg, server, done)

	if err := server.Start(); err != nil {
		if err != http.ErrServerClosed {
			logg.Fatalf("HTTP server: failed to start: %v", err)
		}
	}

	<-done

	if storageSQL, ok := storage.(*sqlstorage.Storage); ok {
		logg.Println("Close sql storage")
		storageSQL.Close(nil)

	}

	logg.Println("Stop service calendar")
	logg.Close()
}

func listenForShutdown(signals chan os.Signal, logg *logger.Logger, server *httpserver.Server, done chan<- struct{}) {
	<-signals
	signal.Stop(signals)

	if err := server.Stop(ctx); err != nil {
		logg.Fatalf("HTTP server: failed to stop: %v", err)
	}

	done <- struct{}{}
}
