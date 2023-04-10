package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Add pure Go Postgres driver for the database/sql package.
	_ "github.com/lib/pq"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/retry"
)

const retriesConnect = 5

type Conf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Storage struct {
	db *sql.DB

	conf *Conf

	insertEvent      *sql.Stmt
	getEvent         *sql.Stmt
	getAllEvents     *sql.Stmt
	getEventsForUser *sql.Stmt

	insertUser  *sql.Stmt
	getUser     *sql.Stmt
	getAllUsers *sql.Stmt
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	var st Storage

	st.conf = conf

	if err := st.RetryConnect(ctx, retriesConnect); err != nil {
		return nil, fmt.Errorf("new database connection: %w", err)
	}

	if err := MigrationsUp(conf); err != nil {
		return nil, fmt.Errorf("database migrations: %w", err)
	}

	if err := st.prepareQuery(ctx); err != nil {
		return nil, fmt.Errorf("prepare query: %w", err)
	}

	return &st, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	dbConfStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.conf.Host, s.conf.Port, s.conf.User, s.conf.Password, s.conf.Name)
	s.db, err = sql.Open("postgres", dbConfStr)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	// Setting db connections pool.
	s.db.SetMaxOpenConns(25)
	s.db.SetMaxIdleConns(25)
	s.db.SetConnMaxLifetime(5 * time.Minute)

	return nil
}

func (s *Storage) RetryConnect(ctx context.Context, retries int) error {
	retryConnect := retry.Retry(s.Connect, retries)
	if err := retryConnect(ctx); err != nil {
		return fmt.Errorf("connection db: %w", err)
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
	s.getAllEvents.Close()
	s.getEventsForUser.Close()

	s.insertUser.Close()
	s.getUser.Close()
	s.getAllUsers.Close()

	return s.db.Close() //nolint:wrapcheck
}
