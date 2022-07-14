package main

import (
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/grpcserver"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/httpserver"
)

type Config struct {
	InMem  bool            `yaml:"inmemory"`
	DB     sqlstorage.Conf `yaml:"db"`
	HTTP   httpserver.Conf `yaml:"http"`
	GRPC   grpcserver.Conf `yaml:"grpc"`
	Logger logger.Conf     `yaml:"logger"`
}
