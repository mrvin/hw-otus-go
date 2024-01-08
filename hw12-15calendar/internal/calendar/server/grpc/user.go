package grpcserver

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-api"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateUser(ctx context.Context, userpb *calendarapi.CreateUserRequest) (*calendarapi.CreateUserResponse, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userpb.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		err := fmt.Errorf("create user: generate hash password: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	user := storage.User{
		Name:         userpb.GetName(),
		HashPassword: string(hashPassword),
		Email:        userpb.GetEmail(),
		Role:         "user",
	}
	id, err := s.authService.CreateUser(ctx, &user)
	if err != nil {
		err = fmt.Errorf("create user: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	uuid := calendarapi.UUID{Value: id.String()}

	return &calendarapi.CreateUserResponse{Id: &uuid}, nil
}

func (s *Server) GetUser(ctx context.Context, _ *calendarapi.GetUserRequest) (*calendarapi.UserResponse, error) {
	userName := GetUserName(ctx)
	if userName == "" {
		err := fmt.Errorf("get user: user name is empty")
		return nil, err
	}
	user, err := s.authService.GetUser(ctx, userName)
	if err != nil {
		err := fmt.Errorf("get user: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	uuid := calendarapi.UUID{Value: user.ID.String()}

	return &calendarapi.UserResponse{
		Id:           &uuid,
		Name:         user.Name,
		HashPassword: user.HashPassword,
		Email:        user.Email,
	}, nil
}

func (s *Server) ListUsers(ctx context.Context, _ *calendarapi.ListUsersRequest) (*calendarapi.ListUsersResponse, error) {
	users, err := s.authService.ListUsers(ctx)
	if err != nil {
		err := fmt.Errorf("get all users: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	pbUsers := make([]*calendarapi.UserResponse, len(users))
	for i, user := range users {
		uuid := calendarapi.UUID{Value: user.ID.String()}
		pbUsers[i] = &calendarapi.UserResponse{
			Id:    &uuid,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return &calendarapi.ListUsersResponse{Users: pbUsers}, nil
}

func (s *Server) UpdateUser(ctx context.Context, userpb *calendarapi.UpdateUserRequest) (*emptypb.Empty, error) {
	userName := GetUserName(ctx)
	if userName == "" {
		err := fmt.Errorf("update user: user name is empty")
		return nil, err
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userpb.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		err := fmt.Errorf("update user: generate hash password: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	user := storage.User{
		Name:         userpb.GetName(),
		HashPassword: string(hashPassword),
		Email:        userpb.GetEmail(),
	}
	if err := s.authService.UpdateUser(ctx, userName, &user); err != nil {
		err := fmt.Errorf("update user: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, _ *calendarapi.DeleteUserRequest) (*emptypb.Empty, error) {
	userName := GetUserName(ctx)
	if userName == "" {
		err := fmt.Errorf("delete user: user name is empty")
		return nil, err
	}
	if err := s.authService.DeleteUser(ctx, userName); err != nil {
		err := fmt.Errorf("delete user: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func GetUserName(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if userName, ok := ctx.Value("username").(string); ok {
		return userName
	}
	return ""
}
