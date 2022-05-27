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

	logFile := logInit(&config.Logger)
	log.Println("Start service calendar")

	var storage app.Storage
	if config.InMem {
		log.Println("Storage in memory")
		storage = memorystorage.New()
	} else {
		log.Println("Storage in sql")
		var storageSQL sqlstorage.Storage
		if err := storageSQL.Connect(ctx, &config.DB); err != nil {
			log.Fatalf("can't connection db: %v", err)
		}
		log.Println("Connect db")
		// if such a tables exists then ignore the error.
		if err := storageSQL.CreateSchemaDB(ctx); err != nil {
			log.Fatalf("can't create schema db: %v", err)
		}
		if err := storageSQL.PrepareQuery(ctx); err != nil {
			log.Fatalf("can't prepare query: %v", err)
		}

		storage = &storageSQL
	}

	server := httpserver.New(&config.HTTP, storage)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)
	done := make(chan struct{})
	go listenForShutdown(signals, server, done)

	if err := server.Start(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalf("HTTP server: failed to start: %v", err)
		}
	}

	<-done

	if storageSQL, ok := storage.(*sqlstorage.Storage); ok {
		log.Println("Close sql storage")
		storageSQL.Close()
	}

	log.Println("Stop service calendar")
	if logFile != nil {
		logFile.Close()
	}
}

func listenForShutdown(signals chan os.Signal, server *httpserver.Server, done chan<- struct{}) {
	<-signals
	signal.Stop(signals)

	if err := server.Stop(ctx); err != nil {
		log.Fatalf("HTTP server: failed to stop: %v", err)
	}

	done <- struct{}{}
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
