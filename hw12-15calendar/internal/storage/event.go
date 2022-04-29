package storage

import (
	"time"
)

type Event struct {
	ID          int
	Title       string
	Description string
	StartTime   time.Time
	StopTime    time.Time
	UserID      int
	//	CreatedAt   time.Time
	//	UpdatedAt   time.Time
}

type User struct {
	ID     int
	Name   string
	Events []*Event
}
