package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func (h *Handler) CreateEvent(res http.ResponseWriter, req *http.Request) {
	event, err := unmarshalEvent(req)
	if err != nil {
		slog.Error("Get event from request body: " + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if event.Title == "" || event.StartTime.IsZero() || event.UserID == 0 {
		errMsg := "Empty fields event: title, start time, user id"
		slog.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	id, err := h.eventService.CreateEvent(req.Context(), event)
	if err != nil {
		slog.Error("Saving event to storage: " + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseID, err := json.Marshal(&responseID{ID: id})
	if err != nil {
		slog.Error("Marshaling response id to json: " + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if _, err := res.Write(jsonResponseID); err != nil {
		slog.Error("Write id to response: " + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

//nolint:dupl
func (h *Handler) GetEvent(res http.ResponseWriter, req *http.Request) {
	id, err := getID(req)
	if err != nil {
		slog.Error("Get event id from request body: " + err.Error())
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	event, err := h.eventService.GetEvent(req.Context(), id)
	if err != nil {
		slog.Error("Get event from storage: " + err.Error())
		if errors.Is(err, storage.ErrNoEvent) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		slog.Error("Marshaling event to json: " + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if _, err := res.Write(jsonEvent); err != nil {
		slog.Error("Write event to response: " + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateEvent(res http.ResponseWriter, req *http.Request) {
	// Update only required fields
	event, err := unmarshalEvent(req)
	if err != nil {
		slog.Error("Get event from request body: " + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if event.ID == 0 {
		errMsg := "Event id not set"
		slog.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := h.eventService.UpdateEvent(req.Context(), event); err != nil {
		slog.Error("Update event in storage: " + err.Error())
		if errors.Is(err, storage.ErrNoEvent) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func (h *Handler) DeleteEvent(res http.ResponseWriter, req *http.Request) {
	id, err := getID(req)
	if err != nil {
		slog.Error("Get event id from request body: " + err.Error())
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := h.eventService.DeleteEvent(req.Context(), id); err != nil {
		slog.Error("Delete event in storage: " + err.Error())
		if errors.Is(err, storage.ErrNoEvent) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func unmarshalEvent(req *http.Request) (*storage.Event, error) {
	var event storage.Event

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read body req: %w", err)
	}

	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("unmarshal body req: %w", err)
	}

	return &event, nil
}
