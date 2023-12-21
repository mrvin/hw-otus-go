package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) (uuid.UUID, error) {
	if err := s.insertUser.QueryRowContext(ctx, user.Name, user.HashPassword, user.Email).Scan(&user.ID); err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	return user.ID, nil
}

func (s *Storage) GetUser(ctx context.Context, name string) (*storage.User, error) {
	var user storage.User

	if err := s.getUser.QueryRowContext(ctx, name).Scan(&user.ID, &user.Name, &user.HashPassword, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", storage.ErrNoUser, name)
		}
		return nil, fmt.Errorf("can't scan user with name: %s: %w", name, err)
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

func (s *Storage) UpdateUser(ctx context.Context, name string, user *storage.User) error {
	res, err := s.db.ExecContext(ctx, "update users set name = $2, email = $3 where name = $1", name, user.Name, user.Email)
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

func (s *Storage) UpdateUserByName(ctx context.Context, user *storage.User) error {
	sqlUpdateUser := `UPDATE users SET hash_password = $2, email = $3 WHERE name = $1`
	res, err := s.db.ExecContext(ctx, sqlUpdateUser, user.Name, user.HashPassword, user.Email)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w: %s", storage.ErrNoUser, user.Name)
	}

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, name string) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE name = $1", name)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w: %s", storage.ErrNoUser, name)
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
		return fmt.Errorf("%w: %s", storage.ErrNoUser, name)
	}

	return nil
}
