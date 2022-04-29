package memorystorage

import (
	"context"
	"sync"

	"github.com/mrvin/hw-otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	// TODO
	mu sync.RWMutex
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id int) (*storage.Event, error) {
	return nil, nil
}

func (s *Storage) GetEventsForUser(ctx context.Context, id int) ([]*storage.Event, error) {
	return nil, nil
}

func (s *Storage) GetListEvents(ctx context.Context) ([]*storage.Event, error) {
	return nil, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	return nil
}

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	return nil
}

func (s *Storage) GetUser(ctx context.Context, id int) (*storage.User, error) {
	return nil, nil
}

func (s *Storage) GetListUsers(ctx context.Context) ([]*storage.User, error) {
	return nil, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *storage.User) error {
	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	return nil
}
