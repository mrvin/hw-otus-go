package queue

import (
	"time"
)

type Conf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type AlertEvent struct {
	EventID     int
	Title       string
	Description string
	StartTime   time.Time
	UserName    string
	UserEmail   string
}
