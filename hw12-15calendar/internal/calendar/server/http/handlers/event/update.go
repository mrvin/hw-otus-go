package handlerevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	httpresponse "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/response"
)

type RequestUpdateEvent struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"       validate:"required,min=2,max=64"`
	Description string    `json:"description" validate:"required,min=2,max=512"`
	StartTime   time.Time `json:"start_time"  validate:"required"`
	StopTime    time.Time `json:"stop_time"   validate:"required"`
}

func (h *Handler) UpdateEvent(res http.ResponseWriter, req *http.Request) {
	// Read json request
	var request RequestUpdateEvent

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		err := fmt.Errorf("UpdateEvent: read body req: %w", err)
		slog.Error(err.Error())
		httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		err := fmt.Errorf("UpdateEvent: unmarshal body request: %w", err)
		slog.Error(err.Error())
		httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Validation
	if err := validator.New().Struct(request); err != nil {
		errors := err.(validator.ValidationErrors)
		err := fmt.Errorf("UpdateEvent: invalid request: %s", errors)
		slog.Error(err.Error())
		httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	event := storage.Event{
		ID:          request.ID,
		Title:       request.Title,
		Description: request.Description,
		StartTime:   request.StartTime,
		StopTime:    request.StopTime,
	}

	if err := h.eventService.UpdateEvent(req.Context(), &event); err != nil {
		err := fmt.Errorf("UpdateEvent: update event in storage: %w", err)
		slog.Error(err.Error())
		if errors.Is(err, storage.ErrNoEvent) {
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		} else {
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Write json response
	httpresponse.WriteOK(res)

	slog.Info("Event information update was successful")
}
