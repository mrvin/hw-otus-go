package handleruser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	httpresponse "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/response"
	"golang.org/x/crypto/bcrypt"
)

type RequestUpdateUser struct {
	UserName string `json:"user_name" validate:"required,min=3,max=20"`
	//	NewUserName string `json:"new_user_name" validate:"omitempty,min=3,max=20"`
	Password string `json:"password" validate:"omitempty,min=6,max=32"`
	Email    string `json:"email" validate:"omitempty,email"`
}

func (h *Handler) UpdateUser(res http.ResponseWriter, req *http.Request) {
	// Read json request
	var request RequestUpdateUser

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		err := fmt.Errorf("UpdateUser: read body request: %w", err.Error())
		slog.Error(err.Error())
		httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
	}

	if err := json.Unmarshal(body, &request); err != nil {
		err := fmt.Errorf("UpdateUser: unmarshal body request: %w", err)
		slog.Error(err.Error())
		httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
	}

	slog.Debug(
		"Update user request",
		slog.String("username", request.UserName),
		//		slog.String("new username", request.NewUserName),
		slog.String("password", request.Password),
		slog.String("email", request.Email),
	)

	if err := validator.New().Struct(request); err != nil {
		errors := err.(validator.ValidationErrors)
		err := fmt.Errorf("UpdateUser: invalid request: %s", errors)
		slog.Error(err.Error())
		httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		err := fmt.Errorf("UpdateUser: generate hash password: %w", err)
		httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	user := storage.User{
		Name:         request.UserName,
		HashPassword: string(hashPassword),
		Email:        request.Email,
	}

	if err := h.authService.UpdateUserByName(req.Context(), &user); err != nil {
		err := fmt.Errorf("UpdateUser: update user in storage: %w", err)
		slog.Error(err.Error())
		if errors.Is(err, storage.ErrNoUserName) {
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		} else {
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Write json response
	httpresponse.WriteOK(res)

	slog.Info("User information update was successful")
}
