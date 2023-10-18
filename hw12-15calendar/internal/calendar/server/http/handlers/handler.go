package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	authservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/auth"
	eventservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/event"
)

var ErrIDEmpty = errors.New("id is empty")

type Handler struct {
	authService  *authservice.AuthService
	eventService *eventservice.EventService
}

func New(auth *authservice.AuthService, events *eventservice.EventService) *Handler {
	return &Handler{
		authService:  auth,
		eventService: events,
	}
}

type responseID struct {
	ID int64 `json:"id"`
}

func getID(req *http.Request) (int64, error) {
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		return 0, ErrIDEmpty
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("convert id: %w", err)
	}

	return id, nil
}
