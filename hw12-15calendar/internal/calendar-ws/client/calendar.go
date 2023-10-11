package client

import (
	"context"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

type Calendar interface {
	CreateUser(ctx context.Context, name, password, email string) (int64, error)
	GetUser(ctx context.Context, id int64) (*storage.User, error)
	UpdateUser(ctx context.Context, name, password, email string) error
	DeleteUser(ctx context.Context, id int64) error

	ListUsers(ctx context.Context) ([]storage.User, error)

	CreateEvent(ctx context.Context,
		title, description string,
		startTime, stopTime time.Time,
		userID int64) (int64, error)
	GetEvent(ctx context.Context, id int64) (*storage.Event, error)
	//	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, id int64) error

	// ListEvents(ctx context.Context) ([]Event, error)
	ListEventsForUser(ctx context.Context, idUser int64, days int) ([]storage.Event, error)
}
