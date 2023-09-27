package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	if err := s.insertUser.QueryRowContext(ctx, user.Name, user.Email).Scan(&user.ID); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *Storage) GetUser(ctx context.Context, id int) (*storage.User, error) {
	var user storage.User

	if err := s.getUser.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %d", storage.ErrNoUser, id)
		}
		return nil, fmt.Errorf("can't scan user with id: %d: %w", id, err)
	}

	// TRANSACTION SQL
	var err error
	user.Events, err = s.ListEventsForUser(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &user, nil
		}
		return nil, fmt.Errorf("can't scan events for user with id: %d: %w", id, err)
	}

	return &user, nil
}

func (s *Storage) ListUsers(ctx context.Context) ([]storage.User, error) {
	users := make([]storage.User, 0)

	rows, err := s.listUsers.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
		}
		return nil, fmt.Errorf("can't get all users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user storage.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return users, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *storage.User) error {
	res, err := s.db.ExecContext(ctx, "update users set name = $2, email = $3 where id = $1", user.ID, user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w: %d", storage.ErrNoUser, user.ID)
	}

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	res, err := s.db.ExecContext(ctx, "delete from users where id = $1", id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w: %d", storage.ErrNoUser, id)
	}

	return nil
}
