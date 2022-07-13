package main

import (
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
	grpcserver "github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/server/grpc"
	httpserver "github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/server/http"
)

type LoggerConf struct {
	FilePath string `yaml:"filepath"`
	Level    string `yaml:"level"`
}

type Config struct {
	InMem  bool            `yaml:"inmemory"`
	DB     sqlstorage.Conf `yaml:"db"`
	HTTP   httpserver.Conf `yaml:"http"`
	GRPC   grpcserver.Conf `yaml:"grpc"`
	Logger LoggerConf      `yaml:"logger"`
}
