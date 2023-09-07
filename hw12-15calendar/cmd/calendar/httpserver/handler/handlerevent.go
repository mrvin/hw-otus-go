package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

// TODO:add return id created event.
func (h *Handler) CreateEvent(res http.ResponseWriter, req *http.Request) {
	event, err := unmarshalEvent(req)
	if err != nil {
		h.log.Errorf("Get event from request body: %v", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if event.Title == "" || event.StartTime.IsZero() || event.UserID == 0 {
		errMsg := "Empty fields event: title, start time, user id"
		h.log.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := h.app.CreateEvent(req.Context(), event); err != nil {
		h.log.Errorf("Saving event to storage: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

//nolint:dupl
func (h *Handler) GetEvent(res http.ResponseWriter, req *http.Request) {
	id, err := getID(req)
	if err != nil {
		h.log.Errorf("Get event id from request body: %v", err)
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	event, err := h.app.GetEvent(req.Context(), id)
	if err != nil {
		h.log.Errorf("Get event from storage: %v", err)
		if errors.Is(err, storage.ErrNoEvent) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		h.log.Errorf("Marshaling event to json: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if _, err := res.Write(jsonEvent); err != nil {
		h.log.Errorf("Write event to response: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateEvent(res http.ResponseWriter, req *http.Request) {
	// Update only required fields
	event, err := unmarshalEvent(req)
	if err != nil {
		h.log.Errorf("Get event from request body: %v", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if event.ID == 0 {
		errMsg := "Event id not set"
		h.log.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := h.app.UpdateEvent(req.Context(), event); err != nil {
		h.log.Errorf("Update event in storage: %v", err)
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
		h.log.Errorf("Get event id from request body: %v", err)
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := h.app.DeleteEvent(req.Context(), id); err != nil {
		h.log.Errorf("Delete event in storage: %v", err)
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
