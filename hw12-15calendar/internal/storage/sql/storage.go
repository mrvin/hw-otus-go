package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
)

type Storage struct {
	db *sql.DB

	insertEvent      *sql.Stmt
	getEvent         *sql.Stmt
	getEventsForUser *sql.Stmt

	insertUser *sql.Stmt
	getUser    *sql.Stmt
}

func (s *Storage) Connect(ctx context.Context, conf *config.DBConf) error {
	var err error
	dbConfStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.User, conf.Password, conf.Name)
	s.db, err = sql.Open("postgres", dbConfStr)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	err = s.db.PingContext(ctx)

	return fmt.Errorf("connection db: %w", err)
}

func (s *Storage) CreateSchemaDB(ctx context.Context) error {
	sqlCreateTableUsers := `
	CREATE TABLE users (
		id serial primary key,
		name text,
		email text
	)`
	_, err := s.db.ExecContext(ctx, sqlCreateTableUsers)
	// if such a table exists then ignore the error.
	if err != nil {
		return fmt.Errorf("create table users: %w", err)
	}

	sqlCreateTableEvents := `
	CREATE TABLE events (
		id serial primary key,
		title text,
		description text,
		start_time timestamptz,
		stop_time timestamptz,
		user_id integer references users(id) on delete cascade
	)`
	_, err = s.db.ExecContext(ctx, sqlCreateTableEvents)
	if err != nil {
		return fmt.Errorf("create table events: %w", err)
	}

	return nil
}

func (s *Storage) PrepareQuery(ctx context.Context) error {
	var err error
	fmtStrErr := "prepare \"%s\" query: %w"
	// Event query prepare
	sqlInsertEvent := "insert into events (title, description, start_time, stop_time, user_id) values ($1, $2, $3, $4, $5) returning id"
	s.insertEvent, err = s.db.PrepareContext(ctx, sqlInsertEvent)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insertEvent", err)
	}
	sqlGetEvent := "select id, title, description, start_time, stop_time, user_id from events where id = $1"
	s.getEvent, err = s.db.PrepareContext(ctx, sqlGetEvent)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "getEvent", err)
	}
	sqlGetEventsForUser := "select id, title, description, start_time, stop_time, user_id from events where user_id = $1"
	s.getEventsForUser, err = s.db.PrepareContext(ctx, sqlGetEventsForUser)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "getEventsForUser", err)
	}

	// User query prepare
	sqlInsertUser := "insert into users (name, email) values ($1, $2) returning id"
	s.insertUser, err = s.db.PrepareContext(ctx, sqlInsertUser)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insertUser", err)
	}
	sqlGetUser := "select id, name, email from users where id = $1"
	s.getUser, err = s.db.PrepareContext(ctx, sqlGetUser)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "getUser", err)
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.insertEvent.Close()
	s.getEvent.Close()
	s.getEventsForUser.Close()

	s.insertUser.Close()
	s.getUser.Close()

	return s.db.Close() //nolint:wrapcheck
}
