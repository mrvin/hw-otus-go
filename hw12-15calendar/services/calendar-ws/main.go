package main

import (
	"errors"
	"flag"
	stdlog "log"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/tracer"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar-ws/grpcclient"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar-ws/httpserver"
)

var infoService = tracer.InfoService{
	Name:    "Calendar-ws",
	Version: "1.0.0",
}

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
	defer log.Sync()

	if err := tracer.TraceInit(&conf.Tracer, &infoService); err != nil {
		log.Errorf("Init jaeger tracer: %v", err)
	}

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
