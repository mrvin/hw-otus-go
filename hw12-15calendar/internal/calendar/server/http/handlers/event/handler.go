package handlerevent

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	eventservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/event"
)

var ErrIDEmpty = errors.New("id is empty")

type Handler struct {
	eventService *eventservice.EventService
}

func New(events *eventservice.EventService) *Handler {
	return &Handler{
		eventService: events,
	}
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
