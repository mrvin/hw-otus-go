package main

import (
	"context"
	"errors"
	"flag"
	stdlog "log"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/metric"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/tracer"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar-ws/grpcclient"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar-ws/httpserver"
)

const serviceName = "Calendar-ws"

var ctx = context.Background()

func main() {
	configFile := flag.String("config", "/etc/calendar/calendar-ws.yml", "path to configuration file")
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

	clientGRPC, err := grpcclient.New(&conf.GRPC)
	if err != nil {
		log.Errorf("gRPC client: %v", err)
		return
	}
	defer clientGRPC.Close()
	log.Info("Connect gRPC server")

	serverHTTP := httpserver.New(&conf.HTTP, clientGRPC.Cl)

	if err := serverHTTP.Start(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("HTTP server failed to start: %v", err)
			return
		}
	}
}
