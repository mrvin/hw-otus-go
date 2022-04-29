package logger

import (
	"log"
	"os"

	"github.com/mrvin/hw-otus-go/hw12_13_14_15_calendar/internal/config"
)

type Logger struct {
	*log.Logger
	logFile *os.File
}

func Create(conf *config.LoggerConf) (*Logger, error) {
	logFile := os.Stdout
	if conf.FilePath != "" {
		var err error
		logFile, err = os.OpenFile(conf.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
	}
	logger := Logger{log.New(logFile, "", log.LstdFlags|log.Lshortfile), logFile}

	return &logger, nil
}

func (l *Logger) Close() {
	l.logFile.Close()
}
