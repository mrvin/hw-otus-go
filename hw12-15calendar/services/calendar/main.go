//go:generate protoc -I=../../api/ --go_out=../../internal/calendarapi --go-grpc_out=require_unimplemented_servers=false:../../internal/calendarapi ../../api/eventservice.proto

package main

import (
	"context"
	"errors"
	"flag"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/metric"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	memorystorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/memory"
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/tracer"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/app"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/grpcserver"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/httpserver"
	"go.uber.org/zap"
)

const serviceName = "Calendar"

// TODO: ctx
var ctx = context.Background()

func main() {
	configFile := flag.String("config", "/etc/calendar/calendar.yml", "path to configuration file")
	flag.Parse()

	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		stdlog.Printf("Parse config: %v", err)
		return
	}

	log, err := logger.LogInit(&conf.Logger)
	if err != nil {
		stdlog.Printf("Init logger: %v", err)
		return
	}
	defer func() {
		if err := log.Sync(); err != nil {
			log.Errorf("logger sync: %v", err)
		}
	}()

	tp, err := tracer.Init(ctx, &conf.Tracer, serviceName)
	if err != nil {
		log.Errorf("Init tracer: %v", err)
		return
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Errorf("Tracer shutdown: %v", err)
		}
	}()

	mp, err := metric.Init(ctx, &conf.Metric, serviceName)
	if err != nil {
		log.Errorf("Init metric : %v", err)
		return
	}
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			log.Errorf("Metric shutdown: %v", err)
		}
	}()

	var storage storage.Storage
	if conf.InMem {
		log.Info("Storage in memory")
		storage = memorystorage.New()
	} else {
		var err error
		log.Info("Storage in sql database")
		storage, err = sqlstorage.New(ctx, &conf.DB)
		if err != nil {
			log.Errorf("SQL database: %v", err)
			return
		}
		log.Info("Connected to database")
	}

	app := app.New(storage)
	serverHTTP := httpserver.New(&conf.HTTP, app)
	serverGRPC, err := grpcserver.New(&conf.GRPC, app)
	if err != nil {
		log.Errorf("New gRPC server: %v", err)
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT /*(Control-C)*/, syscall.SIGTERM)
	go listenForShutdown(signals, serverHTTP, serverGRPC)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		var err error
		if conf.HTTP.IsHTTPS {
			err = serverHTTP.StartTLS(&conf.HTTP.HTTPS)
		} else {
			err = serverHTTP.Start()
		}
		if !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("HTTP server failed to start: %v", err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err := serverGRPC.Start(); err != nil {
			log.Errorf("gRPC server failed to start: %v", err)
			return
		}
	}()

	wg.Wait()

	if storageSQL, ok := storage.(*sqlstorage.Storage); ok {
		if err := storageSQL.Close(); err != nil {
			log.Errorf("Closing the database connection: %v", err)
		} else {
			log.Info("Closing the database connection")
		}
	}

	log.Info("Stop service calendar")
}

func listenForShutdown(signals chan os.Signal, serverHTTP *httpserver.Server, serverGRPC *grpcserver.Server) {
	<-signals
	signal.Stop(signals)

	if err := serverHTTP.Stop(ctx); err != nil {
		zap.S().Errorf("HTTP server failed to stop: %v", err)
		return
	}

	serverGRPC.Stop()
}
