package signup

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	httpresponse "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/response"
	"golang.org/x/crypto/bcrypt"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user *storage.User) error
}

type RequestSignUp struct {
	UserName string `json:"userName" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=32"`
	Email    string `json:"email"    validate:"required,email"`
}

type ResponseSignUp struct {
	Status string `json:"status"`
}

// SignUp registers a new user.
func New(creator UserCreator) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Read json request
		var request RequestSignUp

		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("SignUp: read body request: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("SignUp: unmarshal body request: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		slog.Debug(
			"Sign up request",
			slog.String("username", request.UserName),
			slog.String("password", request.Password),
			slog.String("email", request.Email),
		)

		// Validation
		if err := validator.New().Struct(request); err != nil {
			errors := err.(validator.ValidationErrors)
			err := fmt.Errorf("SignUp: invalid request: %s", errors)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			err := fmt.Errorf("SignUp: generate hash password: %w", err)
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			slog.Error(err.Error())
			return
		}

		user := storage.User{
			Name:         request.UserName,
			HashPassword: string(hashPassword),
			Email:        request.Email,
			Role:         "user",
		}

		if err = creator.CreateUser(req.Context(), &user); err != nil {
			err := fmt.Errorf("SignUp: saving user to storage: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseSignUp{
			Status: "OK",
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("SignUp: marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("SignUp: write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("New user registration was successful")
	}
}
