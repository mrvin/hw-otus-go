package handleruser

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	httpresponse "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/response"
)

type ResponseGetUser struct {
	ID           int64           `json:"id,required"`
	Name         string          `json:"name,required"`
	HashPassword string          `json:"hash_password,required"`
	Email        string          `json:"email,required"`
	Role         string          `json:"role,required"`
	Events       []storage.Event `json:"events,required"`
	Status       string          `json:"status,required"`
}

func (h *Handler) GetUser(res http.ResponseWriter, req *http.Request) {
	userName := handler.GetUserName(req.Context())
	if userName == "" {
		err := fmt.Errorf("GetUser: user name is empty")
		httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.authService.GetUserByName(req.Context(), userName)
	if err != nil {
		err := fmt.Errorf("GetUser: get user from storage: %w", err)
		slog.Error(err.Error())
		if errors.Is(err, storage.ErrNoUserName) {
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		} else {
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Write json response
	response := ResponseGetUser{
		ID:           user.ID,
		Name:         user.Name,
		HashPassword: user.HashPassword,
		Email:        user.Email,
		Role:         user.Role,
		Events:       user.Events,
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
