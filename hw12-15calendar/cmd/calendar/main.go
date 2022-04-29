package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/app"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	httpserver "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/server/http"
	memorystorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/memory"
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
)

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
		if err := storageSQL.Connect(nil, &config.DB); err != nil {
			logg.Fatalf("can't connection db: %v", err)
		}
		storage = &storageSQL
	}

	server := httpserver.New(&config.HTTP, logg, storage)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go listenForShutdown(signals, logg, server)

	if err := server.Start(); err != nil {
		logg.Fatalf("failed to start http server: %v", err)
	}

	logg.Println("Stop service calendar")
}

func listenForShutdown(signals chan os.Signal, logg *logger.Logger, server *httpserver.Server) {
	<-signals
	signal.Stop(signals)

	if err := server.Stop(); err != nil {
		logg.Fatalf("failed to stop http server: %v", err)
	}
	//storage.Close()
	logg.Println("Stop service calendar")
	logg.Close()
}
