package storage

import (
	"context"
	"errors"
	"time"
)

var ErrNoUser = errors.New("no user with id")
var ErrNoEvent = errors.New("no event with id")

type EventStorage interface {
	CreateEvent(ctx context.Context, event *Event) error
	GetEvent(ctx context.Context, id int) (*Event, error)
	GetAllEvents(ctx context.Context) ([]Event, error)
	GetEventsForUser(ctx context.Context, id int) ([]Event, error)
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, id int) error
}

type UserStorage interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int) (*User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) error
}

type Storage interface {
	EventStorage
	UserStorage
}

//nolint:tagliatelle
type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	StopTime    time.Time `json:"stop_time,omitempty"`
	UserID      int       `json:"user_id"`
	//	CreatedAt   time.Time
	//	UpdatedAt   time.Time
}

type User struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Events []Event `json:"events"`
}
