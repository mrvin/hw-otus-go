package grpcclient

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-api"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func (c *Client) Registration(ctx context.Context, name, password, email string) (uuid.UUID, error) {
	user := &calendarapi.CreateUserRequest{
		Name:     name,
		Password: password,
		Email:    email,
	}
	response, err := c.userService.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("gRPC: %w", err)
	}
	slog.Debug("Added user",
		slog.String("id", response.GetId().GetValue()),
		slog.String("username", user.Name),
		slog.String("password", user.Password),
		slog.String("email", user.Email),
	)

	return uuid.MustParse(response.GetId().GetValue()), nil
}

func (c *Client) Login(ctx context.Context, name, password string) (string, error) {
	reqLogin := &calendarapi.LoginRequest{
		Username: name,
		Password: password,
	}
	response, err := c.userService.Login(ctx, reqLogin)
	if err != nil {
		return "", fmt.Errorf("gRPC: %w", err)
	}

	return response.GetAccessToken(), nil
}

func (c *Client) GetUser(ctx context.Context, token string) (*storage.User, error) {
	reqUser := &calendarapi.GetUserRequest{
		AccessToken: token,
	}
	user, err := c.userService.GetUser(ctx, reqUser)
	if err != nil {
		return nil, fmt.Errorf("gRPC: %w", err)
	}

	return &storage.User{
		ID:           uuid.MustParse(user.Id.String()),
		Name:         user.Name,
		HashPassword: user.HashPassword,
		Email:        user.Email,
	}, nil
}

func (c *Client) UpdateUser(ctx context.Context, token, name, password, email string) error {
	user := &calendarapi.UpdateUserRequest{
		AccessToken: token,
		Name:        name,
		Password:    password,
		Email:       email,
	}
	_, err := c.userService.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("gRPC: %w", err)
	}
	slog.Debug("Update user",
		slog.String("username", name),
	)

	return nil
}

func (c *Client) DeleteUser(ctx context.Context, token string) error {
	reqDeleteUser := &calendarapi.DeleteUserRequest{
		AccessToken: token,
	}
	if _, err := c.userService.DeleteUser(ctx, reqDeleteUser); err != nil {
		return fmt.Errorf("gRPC: %w", err)
	}

	return nil
}

func (c *Client) ListUsers(ctx context.Context, token string) ([]storage.User, error) {
	reqListUsers := &calendarapi.ListUsersRequest{
		AccessToken: token,
	}
	users, err := c.userService.ListUsers(ctx, reqListUsers)
	if err != nil {
		return nil, fmt.Errorf("gRPC: %w", err)
	}

	listUsers := make([]storage.User, 0, len(users.Users))
	for _, user := range users.Users {
		listUsers = append(listUsers, storage.User{
			ID:           uuid.MustParse(user.Id.String()),
			Name:         user.Name,
			HashPassword: user.HashPassword,
			Email:        user.Email,
		})
	}

	return listUsers, nil
}
