package main

import (
	"context"
	"errors"
	"flag"
	stdlog "log"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-ws/client"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-ws/client/grpc"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-ws/client/http"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-ws/httpserver"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/metric"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/tracer"
)

const serviceName = "Calendar-ws"

type Config struct {
	Client     string          `yaml:"client"`
	HTTPClient httpclient.Conf `yaml:"http-client"`
	GRPCClient grpcclient.Conf `yaml:"grpc-client"`
	HTTP       httpserver.Conf `yaml:"http"`
	Logger     logger.Conf     `yaml:"logger"`
	Tracer     tracer.Conf     `yaml:"tracer"`
	Metric     metric.Conf     `yaml:"metrics"`
}

var ctx = context.Background()

func main() {
	configFile := flag.String("config", "/etc/calendar/calendar-ws.yml", "path to configuration file")
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
		slog.Info("Init logger")
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

	var client client.Calendar
	if conf.Client == "grpc" {
		client, err = grpcclient.New(ctx, &conf.GRPCClient)
		if err != nil {
			slog.Error("Failed to init gRPC client: " + err.Error())
			return
		}
		slog.Info("Connect to gRPC server")
		defer func() {
			if clientGRPC, ok := client.(*grpcclient.Client); ok {
				if err := clientGRPC.Close(); err != nil {
					slog.Error("Failed to close gRPC connect: " + err.Error())
				} else {
					slog.Info("Closing the gRPC connection")
				}
			}
		}()
	}
	serverHTTP := httpserver.New(&conf.HTTP, client)

	if err := serverHTTP.Start(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start http server: " + err.Error())
			return
		}
	}

	slog.Info("Stop service " + serviceName)
}
