package memorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

type Storage struct {
	mUsers     map[int]storage.User
	maxIDEvent int
	muUsers    sync.RWMutex

	mEvents   map[int]storage.Event
	maxIDUser int
	muEvents  sync.RWMutex
}

func New() *Storage {
	var s Storage
	s.mUsers = make(map[int]storage.User)
	s.mEvents = make(map[int]storage.Event)

	return &s
}

func (s *Storage) CreateUser(_ context.Context, user *storage.User) error {
	s.muUsers.Lock()
	defer s.muUsers.Unlock()

	s.maxIDUser++
	user.ID = s.maxIDUser
	s.mUsers[s.maxIDUser] = *user

	return nil
}

func (s *Storage) GetUser(_ context.Context, id int) (*storage.User, error) {
	s.muUsers.RLock()
	defer s.muUsers.RUnlock()

	user, ok := s.mUsers[id]
	if !ok {
		return nil, fmt.Errorf("%w: %d", storage.ErrNoUser, id)
	}

	user.Events, _ = s.GetEventsForUser(nil, user.ID)

	return &user, nil
}

func (s *Storage) GetAllUsers(_ context.Context) ([]storage.User, error) {
	users := make([]storage.User, 0)

	s.muUsers.RLock()
	for _, user := range s.mUsers {
		users = append(users, user)
	}
	s.muUsers.RUnlock()

	return users, nil
}

func (s *Storage) UpdateUser(_ context.Context, user *storage.User) error {
	s.muUsers.Lock()
	defer s.muUsers.Unlock()

	if _, ok := s.mUsers[user.ID]; !ok {
		return fmt.Errorf("%w: %d", storage.ErrNoUser, user.ID)
	}
	s.mUsers[user.ID] = *user

	return nil
}

func (s *Storage) DeleteUser(_ context.Context, id int) error {
	s.muUsers.Lock()
	if _, ok := s.mUsers[id]; !ok {
		s.muUsers.Unlock()
		return fmt.Errorf("%w: %d", storage.ErrNoUser, id)
	}
	delete(s.mUsers, id)
	s.muUsers.Unlock()

	s.muEvents.Lock()
	for _, event := range s.mEvents {
		if event.UserID == id {
			delete(s.mEvents, event.ID)
		}
	}
	s.muEvents.Unlock()

	return nil
}

func (s *Storage) CreateEvent(_ context.Context, event *storage.Event) error {
	s.muUsers.Lock()
	if _, ok := s.mUsers[event.UserID]; !ok {
		s.muUsers.Unlock()
		return fmt.Errorf("%w: %d", storage.ErrNoUser, event.UserID)
	}
	s.muUsers.Unlock()

	s.muEvents.Lock()
	defer s.muEvents.Unlock()

	s.maxIDEvent++
	event.ID = s.maxIDEvent
	s.mEvents[s.maxIDEvent] = *event

	return nil
}

func (s *Storage) GetEvent(_ context.Context, id int) (*storage.Event, error) {
	s.muEvents.RLock()
	defer s.muEvents.RUnlock()

	user, ok := s.mEvents[id]
	if !ok {
		return nil, fmt.Errorf("%w: %d", storage.ErrNoEvent, id)
	}

	return &user, nil
}

func (s *Storage) GetAllEvents(_ context.Context) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	s.muEvents.RLock()
	for _, event := range s.mEvents {
		events = append(events, event)
	}
	s.muEvents.RUnlock()

	return events, nil
}

func (s *Storage) UpdateEvent(_ context.Context, event *storage.Event) error {
	s.muEvents.Lock()
	defer s.muEvents.Unlock()

	if _, ok := s.mEvents[event.ID]; !ok {
		return fmt.Errorf("%w: %d", storage.ErrNoEvent, event.ID)
	}
	s.mEvents[event.ID] = *event

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id int) error {
	s.muEvents.Lock()
	defer s.muEvents.Unlock()

	if _, ok := s.mEvents[id]; !ok {
		return fmt.Errorf("%w: %d", storage.ErrNoEvent, id)
	}
	delete(s.mEvents, id)

	return nil
}

func (s *Storage) GetEventsForUser(_ context.Context, id int) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	s.muEvents.RLock()
	for _, event := range s.mEvents {
		if event.UserID == id {
			events = append(events, event)
		}
	}
	s.muEvents.RUnlock()

	return events, nil
}
