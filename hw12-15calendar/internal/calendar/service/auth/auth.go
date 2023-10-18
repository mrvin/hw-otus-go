package authservice

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authSt storage.UserStorage
	tr     trace.Tracer // Think about it.
}

func New(st storage.UserStorage) *AuthService {
	return &AuthService{
		st,
		otel.GetTracerProvider().Tracer("Auth service"),
	}
}

func (a *AuthService) CreateUser(ctx context.Context, user *storage.User) (int64, error) {
	cctx, sp := a.tr.Start(ctx, "CreateUser")
	defer sp.End()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("generate hash password: %w", err)
	}
	user.HashPassword = string(hashPassword)

	return a.authSt.CreateUser(cctx, user)
}

func (a *AuthService) GetUser(ctx context.Context, id int64) (*storage.User, error) {
	cctx, sp := a.tr.Start(ctx, "GetUser")
	defer sp.End()

	return a.authSt.GetUser(cctx, id)
}

func (a *AuthService) UpdateUser(ctx context.Context, user *storage.User) error {
	cctx, sp := a.tr.Start(ctx, "UpdateUser")
	defer sp.End()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate hash password: %w", err)
	}
	user.HashPassword = string(hashPassword)

	return a.authSt.UpdateUser(cctx, user)
}

func (a *AuthService) DeleteUser(ctx context.Context, id int64) error {
	cctx, sp := a.tr.Start(ctx, "DeleteUser")
	defer sp.End()

	return a.authSt.DeleteUser(cctx, id)
}
func (a *AuthService) DeleteUserByName(ctx context.Context, name string) error {
	cctx, sp := a.tr.Start(ctx, "DeleteUser")
	defer sp.End()

	return a.authSt.DeleteUserByName(cctx, name)
}

func (a *AuthService) ListUsers(ctx context.Context) ([]storage.User, error) {
	cctx, sp := a.tr.Start(ctx, "GetAllUsers")
	defer sp.End()

	return a.authSt.ListUsers(cctx)
}

const secret = "our-secret"

func (a *AuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	if !a.validCredentials(ctx, username, password) {
		return "", fmt.Errorf("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"iat":      time.Now().Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to create a token")
	}

	return tokenString, nil
}

func (a *AuthService) validCredentials(ctx context.Context, username, password string) bool {
	user, err := a.authSt.GetUserByName(ctx, username)
	if err != nil {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
		return false
	}

	return true
}

// HashPassword returns the bcrypt hash of the password
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}
