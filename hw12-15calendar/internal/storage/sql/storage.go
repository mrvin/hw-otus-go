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
		return err
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) PrepareQuery(ctx context.Context) (err error) {
	// Event query prepare
	sqlInsertEvent := "insert into events (title, description, start_time, stop_time, user_id) values ($1, $2, $3, $4, $5) returning id"
	s.insertEvent, err = s.db.PrepareContext(ctx, sqlInsertEvent)
	if err != nil {
		return
	}
	sqlGetEvent := "select id, title, description, start_time, stop_time, user_id from events where id = $1"
	s.getEvent, err = s.db.PrepareContext(ctx, sqlGetEvent)
	if err != nil {
		return
	}
	sqlGetEventsForUser := "select id, title, description, start_time, stop_time, user_id from events where user_id = $1"
	s.getEventsForUser, err = s.db.PrepareContext(ctx, sqlGetEventsForUser)
	if err != nil {
		return
	}

	// User query prepare
	sqlInsertUser := "insert into users (name, email) values ($1, $2) returning id"
	s.insertUser, err = s.db.PrepareContext(ctx, sqlInsertUser)
	if err != nil {
		return
	}
	sqlGetUser := "select id, name, email from users where id = $1"
	s.getUser, err = s.db.PrepareContext(ctx, sqlGetUser)
	if err != nil {
		return
	}

	return
}

func (s *Storage) Close(ctx context.Context) error {
	s.insertEvent.Close()
	s.getEvent.Close()
	s.getEventsForUser.Close()

	s.insertUser.Close()
	s.getUser.Close()

	return s.db.Close()
}
