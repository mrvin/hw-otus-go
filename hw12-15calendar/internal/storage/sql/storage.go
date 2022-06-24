package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Storage struct {
	db *sql.DB

	insertEvent      *sql.Stmt
	getEvent         *sql.Stmt
	getAllEvents     *sql.Stmt
	getEventsForUser *sql.Stmt

	insertUser  *sql.Stmt
	getUser     *sql.Stmt
	getAllUsers *sql.Stmt
}

func New(ctx context.Context, conf *DBConf) (*Storage, error) {
	var s Storage

	if err := s.connect(ctx, conf); err != nil {
		return nil, fmt.Errorf("connection db: %w", err)
	}
	if err := s.createSchema(ctx); err != nil {
		return nil, fmt.Errorf("create schema db: %w", err)
	}
	if err := s.prepareQuery(ctx); err != nil {
		return nil, fmt.Errorf("prepare query: %w", err)
	}

	return &s, nil
}

func (s *Storage) connect(ctx context.Context, conf *DBConf) error {
	var err error
	dbConfStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Name)
	s.db, err = sql.Open("postgres", dbConfStr)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	return nil
}

func (s *Storage) createSchema(ctx context.Context) error {
	sqlCreateTableUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id serial primary key,
		name text,
		email text
	)`
	if _, err := s.db.ExecContext(ctx, sqlCreateTableUsers); err != nil {
		return fmt.Errorf("create table users: %w", err)
	}

	sqlCreateTableEvents := `
	CREATE TABLE IF NOT EXISTS events (
		id serial primary key,
		title text,
		description text,
		start_time timestamptz,
		stop_time timestamptz,
		user_id integer references users(id) on delete cascade
	)`
	if _, err := s.db.ExecContext(ctx, sqlCreateTableEvents); err != nil {
		return fmt.Errorf("create table events: %w", err)
	}

	return nil
}

func (s *Storage) DropSchemaDB(ctx context.Context) error {
	sqlDropTableEvents := `DROP TABLE IF EXISTS events`
	if _, err := s.db.ExecContext(ctx, sqlDropTableEvents); err != nil {
		return fmt.Errorf("drop table events: %w", err)
	}

	sqlDropTableUsers := `DROP TABLE IF EXISTS users`
	if _, err := s.db.ExecContext(ctx, sqlDropTableUsers); err != nil {
		return fmt.Errorf("drop table users: %w", err)
	}

	return nil
}

func (s *Storage) prepareQuery(ctx context.Context) error {
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
	sqlGetAllEvent := "SELECT * FROM events"
	s.getAllEvents, err = s.db.PrepareContext(ctx, sqlGetAllEvent)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "getAllEvents", err)
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
	sqlGetAllUsers := "SELECT * FROM users"
	s.getAllUsers, err = s.db.PrepareContext(ctx, sqlGetAllUsers)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "getAllUsers", err)
	}

	return nil
}

func (s *Storage) Close() error {
	s.insertEvent.Close()
	s.getEvent.Close()
	s.getEventsForUser.Close()

	s.insertUser.Close()
	s.getUser.Close()

	return s.db.Close() //nolint:wrapcheck
}
