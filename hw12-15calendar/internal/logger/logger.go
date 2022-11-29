package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Conf struct {
	FilePath string `yaml:"filepath"`
	Level    string `yaml:"level"`
}

func LogInit(conf *Conf) error {
	cfg := zap.NewDevelopmentConfig()
	level, err := zap.ParseAtomicLevel(conf.Level)
	if err != nil {
		return fmt.Errorf("parse level: %w", err)
	}
	cfg.Level = level
	cfg.Encoding = "json"

	outputPaths := []string{"stdout"}
	if conf.FilePath != "" {
		outputPaths = append(outputPaths, conf.FilePath)
	}
	cfg.OutputPaths = outputPaths

	cfg.ErrorOutputPaths = []string{"stderr"}

	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.StacktraceKey = "stacktrace"

	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000")

	logger, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("build  logger: %w", err)
	}
	zap.ReplaceGlobals(logger)

	return nil
}
