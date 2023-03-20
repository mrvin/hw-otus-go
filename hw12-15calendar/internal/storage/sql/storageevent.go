package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	if err := s.insertEvent.QueryRowContext(ctx, event.Title, event.Description, event.StartTime, event.StopTime, event.UserID).Scan(&event.ID); err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id int) (*storage.Event, error) {
	var event storage.Event

	if err := s.getEvent.QueryRowContext(ctx, id).Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.StopTime, &event.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: %d", storage.ErrNoEvent, id)
		}
		return nil, fmt.Errorf("can't get event with id: %d: %w", id, err)
	}

	return &event, nil
}

func (s *Storage) GetAllEvents(ctx context.Context) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	rows, err := s.getAllEvents.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return events, nil
		}
		return nil, fmt.Errorf("can't get all events: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var event storage.Event
		err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.StopTime, &event.UserID)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("rows error: %w", err)
	}

	return events, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	res, err := s.db.ExecContext(ctx, "update events set title = $2, description = $3, start_time = $4, stop_time = $5, user_id = $6 where id = $1",
		event.ID, event.Title, event.Description, event.StartTime, event.StopTime, event.UserID)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w: %d", storage.ErrNoEvent, event.ID)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	res, err := s.db.ExecContext(ctx, "delete from events where id = $1", id)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w: %d", storage.ErrNoEvent, id)
	}

	return nil
}

func (s *Storage) GetEventsForUser(ctx context.Context, userID int) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	rows, err := s.getEventsForUser.QueryContext(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return events, nil
		}
		return nil, fmt.Errorf("can't get events for user with id: %d: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var event storage.Event
		err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.StopTime, &event.UserID)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("rows error: %w", err)
	}

	return events, nil
}
