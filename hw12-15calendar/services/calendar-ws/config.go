package main

import (
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar-ws/grpcclient"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar-ws/httpserver"
)

type Config struct {
	HTTP   httpserver.Conf `yaml:"http"`
	GRPC   grpcclient.Conf `yaml:"grpc"`
	Logger logger.Conf     `yaml:"logger"`
}
