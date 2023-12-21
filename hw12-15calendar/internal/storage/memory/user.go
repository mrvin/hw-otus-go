package memorystorage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func (s *Storage) CreateUser(_ context.Context, user *storage.User) (uuid.UUID, error) {
	s.muUsers.Lock()
	defer s.muUsers.Unlock()

	user.ID = uuid.New()
	s.mUsers[user.Name] = *user

	return user.ID, nil
}

func (s *Storage) GetUser(_ context.Context, name string) (*storage.User, error) {
	s.muUsers.RLock()
	defer s.muUsers.RUnlock()

	user, ok := s.mUsers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", storage.ErrNoUser, name)
	}

	return &user, nil
}

func (s *Storage) GetUserByID(_ context.Context, id uuid.UUID) (*storage.User, error) {
	s.muUsers.RLock()
	defer s.muUsers.RUnlock()

	for _, user := range s.mUsers {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("%w: %v", storage.ErrNoUser, id)
}

func (s *Storage) UpdateUser(_ context.Context, name string, user *storage.User) error {
	s.muUsers.Lock()
	defer s.muUsers.Unlock()

	if _, ok := s.mUsers[name]; !ok {
		return fmt.Errorf("%w: %s", storage.ErrNoUser, user.Name)
	}
	if name != user.Name {
		delete(s.mUsers, name)
	}
	s.mUsers[user.Name] = *user

	return nil
}

func (s *Storage) DeleteUser(_ context.Context, name string) error {
	s.muUsers.Lock()
	user, ok := s.mUsers[name]
	if !ok {
		s.muUsers.Unlock()
		return fmt.Errorf("%w: %s", storage.ErrNoUser, name)
	}
	delete(s.mUsers, name)
	s.muUsers.Unlock()

	s.muEvents.Lock()
	for _, event := range s.mEvents {
		if event.UserID == user.ID {
			delete(s.mEvents, event.ID)
		}
	}
	s.muEvents.Unlock()

	return nil
}

func (s *Storage) ListUsers(_ context.Context) ([]storage.User, error) {
	users := make([]storage.User, 0)

	s.muUsers.RLock()
	for _, user := range s.mUsers {
		users = append(users, user)
	}
	s.muUsers.RUnlock()

	return users, nil
}
