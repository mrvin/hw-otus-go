package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

var ErrNoUser = errors.New("no user")

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	return s.insertUser.QueryRowContext(ctx, user.Name, user.Email).Scan(&user.ID)
}

func (s *Storage) GetUser(ctx context.Context, id int) (*storage.User, error) {
	var user storage.User

	if err := s.getUser.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user id: %d: %w", id, ErrNoUser)
		}
		return nil, fmt.Errorf("can't scan user with id: %d: %w", id, err)
	}

	var err error
	user.Events, err = s.GetEventsForUser(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &user, nil
		}
		return nil, fmt.Errorf("can't scan events for user with id: %d: %w", id, err)
	}

	return &user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *storage.User) error {
	_, err := s.db.ExecContext(ctx, "update users set name = $2, email = $3 where id = $1", user.ID, user.Name, user.Email)

	return err
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "delete from users where id = $1", id)

	return err
}
