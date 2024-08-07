package get

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	handler "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	httpresponse "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/response"
)

type UserGetter interface {
	GetUser(ctx context.Context, name string) (*storage.User, error)
}

type ResponseGetUser struct {
	Name         string `json:"name"`
	HashPassword string `json:"hashPassword"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Status       string `json:"status"`
}

func New(getter UserGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userName := handler.GetUserNameFromContext(req.Context())
		if userName == "" {
			err := fmt.Errorf("GetUser: %w", handler.ErrUserNameIsEmpty)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := getter.GetUser(req.Context(), userName)
		if err != nil {
			err := fmt.Errorf("GetUser: get user from storage: %w", err)
			slog.Error(err.Error())
			if errors.Is(err, storage.ErrNoUser) {
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			} else {
				httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Write json response
		response := ResponseGetUser{
			Name:         user.Name,
			HashPassword: user.HashPassword,
			Email:        user.Email,
			Role:         user.Role,
			Status:       "OK",
		}
		jsonResponseUser, err := json.Marshal(response)
		if err != nil {
			err := fmt.Errorf("GetUser: marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponseUser); err != nil {
			err := fmt.Errorf("GetUser: write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("Get user information was successful")
	}
}
