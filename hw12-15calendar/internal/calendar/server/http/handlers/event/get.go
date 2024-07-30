package handlerevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	httpresponse "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/response"
)

//nolint:tagliatelle
type ResponseGetEvent struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	StopTime    time.Time `json:"stop_time,omitempty"`
	UserName    string    `json:"user_name"`
	Status      string    `json:"status"`
}

func (h *Handler) GetEvent(res http.ResponseWriter, req *http.Request) {
	id, err := getID(req)
	if err != nil {
		err := fmt.Errorf("GetEvent: get event id from request url query: %w", err)
		slog.Error(err.Error())
		if errors.Is(err, ErrIDEmpty) {
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		} else {
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	event, err := h.eventService.GetEvent(req.Context(), id)
	if err != nil {
		err := fmt.Errorf("GetEvent: get event from storage: %w", err)
		slog.Error(err.Error())
		if errors.Is(err, storage.ErrNoEvent) {
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		} else {
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Write json response
	response := ResponseGetEvent{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		StopTime:    event.StopTime,
		UserName:    event.UserName,
		Status:      "OK",
	}
	jsonResponseEvent, err := json.Marshal(response)
	if err != nil {
		err := fmt.Errorf("GetEvent: marshal response: %w", err)
		slog.Error(err.Error())
		httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(jsonResponseEvent); err != nil {
		err := fmt.Errorf("GetEvent: write response: %w", err)
		slog.Error(err.Error())
		httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Get event information was successful")
}
