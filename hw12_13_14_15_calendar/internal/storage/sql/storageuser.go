package sqlstorage

import (
	"context"
	"fmt"

	"github.com/mrvin/hw-otus-go/hw12_13_14_15_calendar/internal/storage"
)

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	statement := "insert into users (name) values ($1) returning id"
	stmt, err := s.db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(user.Name).Scan(&user.ID)
}

func (s *Storage) GetUser(ctx context.Context, id int) (*storage.User, error) {
	var user storage.User

	querySQLGetUser := "select id, name from users where id = $1"
	stmt, err := s.db.Prepare(querySQLGetUser)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// TODO: sql: no rows in result set
	if err := stmt.QueryRow(id).Scan(&user.ID, &user.Name); err != nil {
		return nil, fmt.Errorf("can't scan user(id, name): %v", err)
	}

	user.Events, err = s.GetEventsForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//////////////////
func (s *Storage) GetListUsers(ctx context.Context) ([]*storage.User, error) {
	return nil, nil
}

//////////////////

func (s *Storage) UpdateUser(ctx context.Context, user *storage.User) error {
	_, err := s.db.Exec("update users set name = $2 where id = $1", user.ID, user.Name)

	return err
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	_, err := s.db.Exec("delete from users where id = $1", id)

	return err
}
