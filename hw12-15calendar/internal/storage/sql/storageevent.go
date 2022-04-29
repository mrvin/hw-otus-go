package sqlstorage

import (
	"context"
	"fmt"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	statement := "insert into events (title, description, start_time, stop_time, user_id) values ($1, $2, $3, $4, $5) returning id"
	stmt, err := s.db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(event.Title, event.Description, event.StartTime, event.StopTime, event.UserID).Scan(&event.ID)
}

func (s *Storage) GetEvent(ctx context.Context, id int) (*storage.Event, error) {
	var event storage.Event
	statement := "select id, title, description, start_time, stop_time, user_id from events where id = $1"
	stmt, err := s.db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.QueryRow(id).Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.StopTime, &event.UserID); err != nil {
		return nil, err
	}

	return &event, nil
}

/////////////////
func (s *Storage) GetEventsForUser(ctx context.Context, id int) ([]*storage.Event, error) {
	querySQLGetEventsForUser := "select id, title, description, start_time, stop_time, user_id from events where user_id = $1"
	stmt, err := s.db.Prepare(querySQLGetEventsForUser)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	events := make([]*storage.Event, 0) //TAI
	for rows.Next() {
		var event storage.Event
		err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.StopTime, &event.UserID)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %v", err)
		}
		events = append(events, &event)
	}
	rows.Close()

	return events, nil
}

func (s *Storage) GetListEvents(ctx context.Context) ([]*storage.Event, error) {
	return nil, nil
}

////////////////
func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	_, err := s.db.Exec("update events set title = $2, description = $3, start_time = $4, stop_time = $5, user_id = $6 where id = $1", event.ID, event.Title, event.Description, event.StartTime, event.StopTime, event.UserID)

	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	_, err := s.db.Exec("delete from events where id = $1", id)

	return err
}
