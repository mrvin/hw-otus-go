package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) (int64, error) {
	var id int64
	if err := s.insertUser.QueryRowContext(ctx, user.Name, user.HashPassword, user.Email).Scan(&id); err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, id int64) (*storage.User, error) {
	var user storage.User

	if err := s.getUser.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.HashPassword, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %d", storage.ErrNoUser, id)
		}
		return nil, fmt.Errorf("can't scan user with id: %d: %w", id, err)
	}

	// TODO: TRANSACTION SQL
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

func (s *Storage) GetUserByName(ctx context.Context, name string) (*storage.User, error) {
	var user storage.User

	if err := s.getUserByName.QueryRowContext(ctx, name).Scan(&user.ID, &user.Name, &user.HashPassword, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", storage.ErrNoUserName, name)
		}
		return nil, fmt.Errorf("can't scan user with name: %s: %w", name, err)
	}

	// TODO: TRANSACTION SQL
	var err error
	user.Events, err = s.ListEventsForUser(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &user, nil
		}
		return nil, fmt.Errorf("can't scan events for user with name: %s: %w", name, err)
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
		err = rows.Scan(&user.ID, &user.Name, &user.HashPassword, &user.Email)
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

func (s *Storage) DeleteUser(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
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

func (s *Storage) DeleteUserByName(ctx context.Context, name string) error {
	res, err := s.db.ExecContext(ctx, "delete from users where name = $1", name)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w: %s", storage.ErrNoUserName, name)
	}

	return nil
}
