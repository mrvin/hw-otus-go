package storage

import (
	"errors"
	"time"
)

var ErrNoUser = errors.New("no user with id")
var ErrNoEvent = errors.New("no event with id")

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	StopTime    time.Time `json:"stopTime"`
	UserID      int       `json:"userID"` //nolint:tagliatelle
	//	CreatedAt   time.Time
	//	UpdatedAt   time.Time
}

type User struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Events []Event `json:"events"`
}
