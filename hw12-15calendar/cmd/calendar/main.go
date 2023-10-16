//go:generate protoc -I=../../api/ --go_out=../../internal/calendar-api --go-grpc_out=require_unimplemented_servers=false:../../internal/calendar-api ../../api/event_service.proto
//go:generate protoc -I=../../api/ --go_out=../../internal/calendar-api --go-grpc_out=require_unimplemented_servers=false:../../internal/calendar-api ../../api/user_service.proto
package main

import (
	"context"
	"errors"
	"flag"
	stdlog "log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/app"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/grpcserver"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/httpserver"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/metric"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	memorystorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/memory"
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/tracer"
)

const serviceName = "Calendar"

type Config struct {
	InMem  bool            `yaml:"inmemory"`
	DB     sqlstorage.Conf `yaml:"db"`
	HTTP   httpserver.Conf `yaml:"http"`
	GRPC   grpcserver.Conf `yaml:"grpc"`
	Logger logger.Conf     `yaml:"logger"`
	Tracer tracer.Conf     `yaml:"tracer"`
	Metric metric.Conf     `yaml:"metrics"`
}

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

	logFile, err := logger.Init(&conf.Logger)
	if err != nil {
		stdlog.Printf("Init logger: %v\n", err)
		return
	} else {
		slog.Info("Init logger", slog.String("Logging level", conf.Logger.Level))
		defer func() {
			if err := logFile.Close(); err != nil {
				slog.Error("Close log file: " + err.Error())
			}
		}()
	}
	if conf.Tracer.Enable {
		ctxTracer, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		tp, err := tracer.Init(ctxTracer, &conf.Tracer, serviceName)
		if err != nil {
			slog.Warn("Failed to init tracer: " + err.Error())
		} else {
			slog.Info("Init tracer")
			defer func() {
				if err := tp.Shutdown(ctx); err != nil {
					slog.Error("Failed to shutdown tracer: " + err.Error())
				}
			}()
		}
	}

	if conf.Metric.Enable {
		ctxMetric, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		mp, err := metric.Init(ctxMetric, &conf.Metric, serviceName)
		if err != nil {
			slog.Warn("Failed to init metric: " + err.Error())
		} else {
			slog.Info("Init metric")
			defer func() {
				if err := mp.Shutdown(ctx); err != nil {
					slog.Error("Failed to shutdown metric: " + err.Error())
				}
			}()
		}
	}

	var storage storage.Storage
	if conf.InMem {
		slog.Info("Storage in memory")
		storage = memorystorage.New()
	} else {
		var err error
		slog.Info("Storage in sql database")
		storage, err = sqlstorage.New(ctx, &conf.DB)
		if err != nil {
			slog.Error("Failed to init storage: " + err.Error())
			return
		}
		slog.Info("Connected to database")
	}

	app := app.New(storage)
	serverHTTP := httpserver.New(&conf.HTTP, app)
	serverGRPC, err := grpcserver.New(&conf.GRPC, app)
	if err != nil {
		slog.Error("Failed to init gRPC server: " + err.Error())
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
			slog.Error("Failed to start http server: " + err.Error())
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err := serverGRPC.Start(); err != nil {
			slog.Error("Failed to start gRPC server: " + err.Error())
			return
		}
	}()

	wg.Wait()

	if storageSQL, ok := storage.(*sqlstorage.Storage); ok {
		if err := storageSQL.Close(); err != nil {
			slog.Error("Failed to close storage: " + err.Error())
		} else {
			slog.Info("Closing the database connection")
		}
	}

	slog.Info("Stop service " + serviceName)
}

func listenForShutdown(signals chan os.Signal, serverHTTP *httpserver.Server, serverGRPC *grpcserver.Server) {
	<-signals
	signal.Stop(signals)

	if err := serverHTTP.Stop(ctx); err != nil {
		slog.Error("Failed to stop http server: " + err.Error())
		return
	}

	serverGRPC.Stop()
}
