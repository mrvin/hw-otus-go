package main

import (
	"context"
	"flag"
	stdlog "log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/queue"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/queue/rabbitmq"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/sender/email"
)

type Config struct {
	Queue  queue.Conf  `yaml:"queue"`
	Email  email.Conf  `yaml:"email"`
	Logger logger.Conf `yaml:"logger"`
}

func main() {
	configFile := flag.String("config", "/etc/calendar/sender.yml", "path to configuration file")
	flag.Parse()

	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		stdlog.Printf("Parse config: %v", err)
		return
	}

	log, err := logger.Init(&conf.Logger)
	if err != nil {
		stdlog.Printf("Init logger: %v", err)
		return
	}
	defer log.Sync()

	var qm rabbitmq.Queue

	url := rabbitmq.QueryBuildAMQP(&conf.Queue)

	if err := qm.ConnectAndCreate(url, conf.Queue.Name); err != nil {
		log.Errorf("New queue connection: %v", err)
		return
	}
	defer qm.Close()
	log.Info("Сonnected to queue")

	chConsume, err := qm.GetConsumeChan()
	if err != nil {
		log.Error(err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT /*(Control-C)*/, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	for {
		select {
		case msg, ok := <-chConsume:
			if !ok {
				return
			}
			alertEvent, err := queue.DecodeAlertEvent(msg.Body)
			if err != nil {
				log.Error(err)
				continue
			}

			log.Infof("Take alert message from queue with id: %d\n", alertEvent.EventID)
			emailMsg := email.Message{To: alertEvent.UserEmail, Subject: alertEvent.Title, Description: alertEvent.Description}
			sendEvent(&conf.Email, &emailMsg)
		case <-ctx.Done():
			log.Info("Stop sender")
			return
		}
	}
}

func sendEvent(conf *email.Conf, msg *email.Message) {
	if err := email.Alert(conf, msg); err != nil {
		zap.S().Error(err)
		return
	}
	zap.S().Infof("'%s' event notification sent to '%s'\n", msg.Subject, msg.To)
}
