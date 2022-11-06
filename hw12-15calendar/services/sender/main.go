package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"

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
		fmt.Println(err)
		return
	}
	defer qm.Close()
	log.Println("Сonnected to queue")

	chConsume, err := qm.GetConsumeChan()
	if err != nil {
		fmt.Println(err)
		return
	}

	for msg := range chConsume {
		var eventMsg queue.AlertEvent
		buffer := bytes.NewBuffer(msg.Body)
		dec := gob.NewDecoder(buffer)
		if err := dec.Decode(&eventMsg); err != nil {
			fmt.Println(err)
			continue
		}
		log.Printf("Take alert message from queue with id: %d\n", eventMsg.EventID)
		msg := email.Message{To: []string{eventMsg.UserEmail}, Subject: eventMsg.Title, Description: eventMsg.Description}
		sendEvent(&conf.Email, &msg)
	}
}

func sendEvent(conf *email.Conf, msg *email.Message) {
	if err := email.Alert(conf, msg); err != nil {
		log.Print(err)
		return
	}
	fmt.Printf("'%s' event notification sent to '%s'\n", msg.Subject, msg.To[0])
}
