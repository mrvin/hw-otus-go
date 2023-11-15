package authservice

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

//nolint:tagliatelle
type Conf struct {
	SecretKey           string `yaml:"secret_key"`
	TokenValidityPeriod int    `yaml:"token_validity_period"`
}

type AuthService struct {
	conf   *Conf
	authSt storage.UserStorage
	tr     trace.Tracer // Think about it.
}

func New(st storage.UserStorage, conf *Conf) *AuthService {
	return &AuthService{
		conf,
		st,
		otel.GetTracerProvider().Tracer("Auth service"),
	}
}

func (a *AuthService) CreateUser(ctx context.Context, user *storage.User) (int64, error) {
	cctx, sp := a.tr.Start(ctx, "CreateUser")
	defer sp.End()

	slog.Debug(
		"Create user",
		slog.String("Name", user.Name),
		slog.String("Hash Password", user.HashPassword),
		slog.String("Email", user.Email),
		slog.String("Role", user.Role),
	)

	return a.authSt.CreateUser(cctx, user)
}

func (a *AuthService) GetUserByID(ctx context.Context, id int64) (*storage.User, error) {
	cctx, sp := a.tr.Start(ctx, "GetUser")
	defer sp.End()

	return a.authSt.GetUser(cctx, id)
}

func (a *AuthService) GetUserByName(ctx context.Context, userName string) (*storage.User, error) {
	cctx, sp := a.tr.Start(ctx, "GetUser")
	defer sp.End()

	return a.authSt.GetUserByName(cctx, userName)
}

func (a *AuthService) UpdateUserByName(ctx context.Context, user *storage.User) error {
	cctx, sp := a.tr.Start(ctx, "UpdateUserByName")
	defer sp.End()

	return a.authSt.UpdateUserByName(cctx, user)
}

func (a *AuthService) UpdateUser(ctx context.Context, user *storage.User) error {
	cctx, sp := a.tr.Start(ctx, "UpdateUser")
	defer sp.End()

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

func (a *AuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	if err := a.validCredentials(ctx, username, password); err != nil {
		return "", fmt.Errorf("invalid credentials: %w", err)
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"iat":      time.Now().Unix(),                                                              // IssuedAt
			"exp":      time.Now().Add(time.Duration(a.conf.TokenValidityPeriod) * time.Minute).Unix(), // ExpiresAt
		},
	)
	tokenString, err := token.SignedString([]byte(a.conf.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to create a token: %w", err)
	}

	return tokenString, nil
}

func (a *AuthService) validCredentials(ctx context.Context, username, password string) error {
	user, err := a.authSt.GetUserByName(ctx, username)
	if err != nil {
		return fmt.Errorf("get user by name: %w", err)
	}

	slog.Debug(
		"Compare hash and password",
		slog.String("password", password),
		slog.String("hashedPassword", user.HashPassword),
	)
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
		return fmt.Errorf("compare hash and password: %w", err)
	}

	return nil
}

// Authorized is middleware
func (a *AuthService) Authorized(next http.HandlerFunc) http.HandlerFunc {
	handler := func(res http.ResponseWriter, req *http.Request) {
		authHeaderValue := req.Header.Get("Authorization")
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeaderValue, bearerPrefix) {
			http.Error(res, "request does not contain an Authorization Bearer token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeaderValue, bearerPrefix)

		// Validate token.
		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(a.conf.SecretKey), nil
			},
		)
		if err != nil {
			err := fmt.Errorf("invalid token: %w", err)
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(res, "invalid token", http.StatusUnauthorized)
			return
		}

		username := claims["username"]
		ctx := context.WithValue(req.Context(), "username", username)

		next(res, req.WithContext(ctx)) // Pass request to next handler
	}

	return http.HandlerFunc(handler)
}
