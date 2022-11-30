package main

import (
	"context"
	"flag"
	"fmt"
	stdlog "log"
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
	SchedPeriod int             `yaml:"schedule_period"`
}

func main() {
	configFile := flag.String("config", "/etc/calendar/scheduler.yaml", "path to configuration file")
	flag.Parse()

	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		stdlog.Printf("Parse config: %v", err)
		return
	}

	log, err := logger.LogInit(&conf.Logger)
	if err != nil {
		stdlog.Printf("Init logger: %v", err)
		return
	}
	defer log.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	st, err := sqlstorage.New(ctx, &conf.DB)
	if err != nil {
		log.Errorf("New database connection: %v", err)
		return
	}
	defer st.Close()
	log.Info("Connected to database")

	var qm rabbitmq.Queue

	url := rabbitmq.QueryBuildAMQP(&conf.Queue)

	if err := qm.ConnectAndCreate(url, conf.Queue.Name); err != nil {
		log.Errorf("New queue connection: %v", err)
		return
	}
	defer qm.Close()
	log.Info("Connected to queue")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT /*(Control-C)*/, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	fmt.Println(conf.SchedPeriod)
	schedPeriod := time.Duration(conf.SchedPeriod) * time.Minute
	ticker := time.Tick(schedPeriod)
	for {
		ctxGetAllEvents, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		events, err := st.GetAllEvents(ctxGetAllEvents)
		if err != nil {
			log.Errorf("List event: %v", err)
		}
		log.Info("Start send event...")
		for _, event := range events {
			if cancelled(ctx) {
				break
			}
			nowTime := time.Now()
			if event.StartTime.After(nowTime) && event.StartTime.Before(nowTime.Add(schedPeriod)) {
				user, err := st.GetUser(ctx, event.UserID)
				if err != nil {
					log.Error(err)
					continue
				}
				alertEvent := queue.AlertEvent{EventID: event.ID, Title: event.Title, Description: event.Description,
					StartTime: event.StartTime, UserName: user.Name, UserEmail: user.Email}

				byteAlertEvent, err := queue.EncodeAlertEvent(&alertEvent)
				if err != nil {
					log.Error(err)
					continue
				}

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				if err := qm.SendMsg(ctx, byteAlertEvent); err != nil {
					log.Error(err)
					continue
				}
				log.Infof("Put alert message in queue with id: %d\n", event.ID)
			}
		}
		select {
		case <-ticker:
		// do nothing.
		case <-ctx.Done():
			log.Info("Stop scheduler")
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
