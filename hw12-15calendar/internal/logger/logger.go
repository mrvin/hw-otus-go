package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Conf struct {
	FilePath string `yaml:"filepath"`
	Level    string `yaml:"level"`
}

func Init(conf *Conf) (*os.File, error) {
	var level slog.Level
	level.UnmarshalText([]byte(conf.Level))

	var err error
	var logFile *os.File
	if conf.FilePath == "" {
		logFile = os.Stdout
	} else {
		logFile, err = os.OpenFile(conf.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed open log file: %w", err)
		}
	}
	// todo: setting format cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000")
	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logFile, nil
}
