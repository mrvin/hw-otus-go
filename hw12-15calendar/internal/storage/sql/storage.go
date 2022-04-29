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
}

func (s *Storage) Connect(ctx context.Context, conf *config.DBConf) error {
	var err error
	dbConfStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.User, conf.Password, conf.Name)
	s.db, err = sql.Open("postgres", dbConfStr)
	if err != nil {
		return err
	}

	return s.db.Ping()
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}
