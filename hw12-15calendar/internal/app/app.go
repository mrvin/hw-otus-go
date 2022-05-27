package app

import (
	"context"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

type Storage interface {
	CreateEvent(ctx context.Context, event *storage.Event) error
	GetEvent(ctx context.Context, id int) (*storage.Event, error)
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, id int) error

	GetEventsForUser(ctx context.Context, id int) ([]storage.Event, error)

	CreateUser(ctx context.Context, user *storage.User) error
	GetUser(ctx context.Context, id int) (*storage.User, error)
	UpdateUser(ctx context.Context, user *storage.User) error
	DeleteUser(ctx context.Context, id int) error
}
