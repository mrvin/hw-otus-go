package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/queue"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/queue/rabbitmq"
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
)

type Config struct {
	Queue  queue.Conf      `yaml:"queue"`
	DB     sqlstorage.Conf `yaml:"db"`
	Logger logger.Conf     `yaml:"logger"`
}

const periodSched = time.Minute * 2

func main() {
	configFile := flag.String("config", "/etc/calendar/scheduler.yaml", "path to configuration file")
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

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	st, err := sqlstorage.New(ctx, &conf.DB)
	if err != nil {
		log.Printf("db: %v", err)
		return
	}
	defer st.Close()
	log.Println("Сonnected to database")

	var qm rabbitmq.Queue

	url := rabbitmq.QueryBuildAMQP(&conf.Queue)

	if err := qm.ConnectAndCreate(url, conf.Queue.Name); err != nil {
		fmt.Println(err)
		return
	}
	defer qm.Close()
	log.Println("Сonnected to queue")

	for {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		events, err := st.GetAllEvents(ctx)
		if err != nil {
			fmt.Println(err)
		}
		log.Println("Start send event...")
		for _, event := range events {
			nowTime := time.Now()
			if event.StartTime.After(nowTime) && event.StartTime.Before(nowTime.Add(periodSched)) {
				user, err := st.GetUser(ctx, event.UserID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				eventMsg := queue.AlertEvent{EventID: event.ID, Title: event.Title, Description: event.Description,
					StartTime: event.StartTime, UserName: user.Name, UserEmail: user.Email}

				buffer := new(bytes.Buffer)
				encoder := gob.NewEncoder(buffer)
				if err := encoder.Encode(eventMsg); err != nil {
					fmt.Println(err)
					continue
				}

				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				if err := qm.SendMsg(ctx, buffer.Bytes()); err != nil {
					fmt.Println(err)
					continue
				}
				log.Printf("Put alert message in queue with id: %d\n", event.ID)
			}
		}
		time.Sleep(periodSched)
	}
}
