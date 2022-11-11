package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/queue"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/queue/rabbitmq"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/sender/email"
)

type Config struct {
	Queue  queue.Conf  `yaml:"queue"`
	Email  email.Conf  `yaml:"email"`
	Logger logger.Conf `yaml:"logger"`
}

func main() {
	configFile := flag.String("config", "/etc/calendar/sender.yaml", "path to configuration file")
	flag.Parse()

	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		log.Printf("Parse config: %v", err)
		return
	}

	logFile := logger.LogInit(&conf.Logger)
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	var qm rabbitmq.Queue

	url := rabbitmq.QueryBuildAMQP(&conf.Queue)

	if err := qm.ConnectAndCreate(url, conf.Queue.Name); err != nil {
		log.Println(err)
		return
	}
	defer qm.Close()
	log.Println("Сonnected to queue")

	chConsume, err := qm.GetConsumeChan()
	if err != nil {
		log.Println(err)
		return
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT /*(Control-C)*/, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case msg, ok := <-chConsume:
			if !ok {
				return
			}
			alertEvent, err := queue.DecodeAlertEvent(msg.Body)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("Take alert message from queue with id: %d\n", alertEvent.EventID)
			emailMsg := email.Message{To: alertEvent.UserEmail, Subject: alertEvent.Title, Description: alertEvent.Description}
			sendEvent(&conf.Email, &emailMsg)
		case <-ctx.Done():
			log.Println("Stop sender")
			return
		}
	}
}

func sendEvent(conf *email.Conf, msg *email.Message) {
	if err := email.Alert(conf, msg); err != nil {
		log.Print(err)
		return
	}
	log.Printf("'%s' event notification sent to '%s'\n", msg.Subject, msg.To)
}
