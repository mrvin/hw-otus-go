package memorystorage

import (
	"context"
	"fmt"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) (int64, error) {
	if _, err := s.GetUserByID(ctx, event.UserID); err != nil {
		return 0, err
	}

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

func (s *Storage) ListEventsForUser(ctx context.Context, name string) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	user, err := s.GetUser(ctx, name)
	if err != nil {
		return nil, err
	}

	s.muEvents.RLock()
	for _, event := range s.mEvents {
		if event.UserID == user.ID {
			events = append(events, event)
		}
	}
	s.muEvents.RUnlock()

	return events, nil
}
