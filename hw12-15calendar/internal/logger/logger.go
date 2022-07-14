package logger

import (
	"log"
	"os"
)

type Conf struct {
	FilePath string `yaml:"filepath"`
	Level    string `yaml:"level"`
}

func LogInit(conf *Conf) *os.File {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if conf.FilePath == "" {
		return nil
	}

	logFile, err := os.OpenFile(conf.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("log init: %v", err)
		return nil
	}
	log.SetOutput(logFile)

	return logFile
}
