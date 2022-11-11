package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/logger"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/queue"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/queue/rabbitmq"
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
)

type Config struct {
	Queue       queue.Conf      `yaml:"queue"`
	DB          sqlstorage.Conf `yaml:"db"`
	Logger      logger.Conf     `yaml:"logger"`
	PeriodSched time.Duration   `yaml:"schedule_period"`
}

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
		log.Println(err)
		return
	}
	defer qm.Close()
	log.Println("Сonnected to queue")

	ctx, _ = signal.NotifyContext(context.Background(), syscall.SIGINT /*(Control-C)*/, syscall.SIGTERM, syscall.SIGQUIT)
	//FIXIT
	if conf.PeriodSched == 0 {
		conf.PeriodSched = 2
	}
	periodSched := conf.PeriodSched * time.Minute
	ticker := time.Tick(periodSched)
	for {
		ctxGetAllEvents, _ := context.WithTimeout(context.Background(), 5*time.Second)

		events, err := st.GetAllEvents(ctxGetAllEvents)
		if err != nil {
			log.Println(err)
		}
		log.Println("Start send event...")
		for _, event := range events {
			if cancelled(ctx) {
				break
			}
			nowTime := time.Now()
			if event.StartTime.After(nowTime) && event.StartTime.Before(nowTime.Add(periodSched)) {
				user, err := st.GetUser(ctx, event.UserID)
				if err != nil {
					log.Println(err)
					continue
				}
				alertEvent := queue.AlertEvent{EventID: event.ID, Title: event.Title, Description: event.Description,
					StartTime: event.StartTime, UserName: user.Name, UserEmail: user.Email}

				byteAlertEvent, err := queue.EncodeAlertEvent(&alertEvent)
				if err != nil {
					log.Println(err)
					continue
				}

				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				if err := qm.SendMsg(ctx, byteAlertEvent); err != nil {
					log.Println(err)
					continue
				}
				log.Printf("Put alert message in queue with id: %d\n", event.ID)
			}
		}
		select {
		case <-ticker:
		// do nothing.
		case <-ctx.Done():
			log.Println("Stop scheduler")
			return
		}
	}
}

func cancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
