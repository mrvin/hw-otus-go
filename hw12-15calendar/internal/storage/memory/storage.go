package memorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

type Storage struct {
	mUsers     map[int64]storage.User
	maxIDEvent int64
	muUsers    sync.RWMutex

	mEvents   map[int64]storage.Event
	maxIDUser int64
	muEvents  sync.RWMutex
}

func New() *Storage {
	var s Storage
	s.mUsers = make(map[int64]storage.User)
	s.mEvents = make(map[int64]storage.Event)

	return &s
}

func (s *Storage) CreateUser(_ context.Context, user *storage.User) (int64, error) {
	s.muUsers.Lock()
	defer s.muUsers.Unlock()

	s.maxIDUser++
	user.ID = s.maxIDUser
	s.mUsers[s.maxIDUser] = *user

	return user.ID, nil
}

func (s *Storage) GetUser(_ context.Context, id int64) (*storage.User, error) {
	s.muUsers.RLock()
	defer s.muUsers.RUnlock()

	user, ok := s.mUsers[id]
	if !ok {
		return nil, fmt.Errorf("%w: %d", storage.ErrNoUser, id)
	}
	//nolint:contextcheck
	user.Events, _ = s.ListEventsForUser(context.TODO(), user.ID)

	return &user, nil
}

func (s *Storage) GetUserByName(_ context.Context, name string) (*storage.User, error) {
	var resUser storage.User

	s.muUsers.RLock()
	defer s.muUsers.RUnlock()

	for _, user := range s.mUsers {
		if user.Name == name {
			resUser = user
			break
		}
	}

	if resUser.ID == 0 {
		return nil, fmt.Errorf("%w: %s", storage.ErrNoUserName, name)
	}
	//nolint:contextcheck
	resUser.Events, _ = s.ListEventsForUser(context.TODO(), resUser.ID)

	return &resUser, nil
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

func (s *Storage) UpdateUser(_ context.Context, user *storage.User) error {
	s.muUsers.Lock()
	defer s.muUsers.Unlock()

	if _, ok := s.mUsers[user.ID]; !ok {
		return fmt.Errorf("%w: %d", storage.ErrNoUser, user.ID)
	}
	s.mUsers[user.ID] = *user

	return nil
}

func (s *Storage) UpdateUserByName(_ context.Context, user *storage.User) error {
	flag := false

	s.muUsers.Lock()
	for id, mUser := range s.mUsers {
		if mUser.Name == user.Name {
			mUser.HashPassword = user.HashPassword
			mUser.Email = user.Email
			s.mUsers[id] = mUser
			flag = true
			break
		}
	}
	s.muUsers.Unlock()

	if !flag {
		return fmt.Errorf("%w: %s", storage.ErrNoUserName, user.Name)
	}
	return nil
}

func (s *Storage) DeleteUser(_ context.Context, id int64) error {
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

func (s *Storage) DeleteUserByName(_ context.Context, name string) error {
	var id int64
	s.muUsers.Lock()
	for _, user := range s.mUsers {
		if user.Name == name {
			delete(s.mUsers, user.ID)
			id = user.ID
			break
		}
	}
	s.muUsers.Unlock()
	if id == 0 {
		return fmt.Errorf("%w: %s", storage.ErrNoUserName, name)
	}

	s.muEvents.Lock()
	for _, event := range s.mEvents {
		if event.UserID == id {
			delete(s.mEvents, event.ID)
		}
	}
	s.muEvents.Unlock()

	return nil
}

func (s *Storage) CreateEvent(_ context.Context, event *storage.Event) (int64, error) {
	s.muUsers.Lock()
	if _, ok := s.mUsers[event.UserID]; !ok {
		s.muUsers.Unlock()
		return 0, fmt.Errorf("%w: %d", storage.ErrNoUser, event.UserID)
	}
	s.muUsers.Unlock()

	s.muEvents.Lock()
	defer s.muEvents.Unlock()

	s.maxIDEvent++
	event.ID = s.maxIDEvent
	s.mEvents[s.maxIDEvent] = *event

	return event.ID, nil
}

func (s *Storage) GetEvent(_ context.Context, id int64) (*storage.Event, error) {
	s.muEvents.RLock()
	defer s.muEvents.RUnlock()

	user, ok := s.mEvents[id]
	if !ok {
		return nil, fmt.Errorf("%w: %d", storage.ErrNoEvent, id)
	}

	return &user, nil
}

func (s *Storage) ListEvents(_ context.Context) ([]storage.Event, error) {
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

func (s *Storage) DeleteEvent(_ context.Context, id int64) error {
	s.muEvents.Lock()
	defer s.muEvents.Unlock()

	if _, ok := s.mEvents[id]; !ok {
		return fmt.Errorf("%w: %d", storage.ErrNoEvent, id)
	}
	delete(s.mEvents, id)

	return nil
}

func (s *Storage) ListEventsForUser(_ context.Context, id int64) ([]storage.Event, error) {
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
