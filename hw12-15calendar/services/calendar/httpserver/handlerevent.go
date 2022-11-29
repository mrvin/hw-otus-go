package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

// TODO:add return id created event
func handleCreateEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	event, err := unmarshalEvent(req)
	if err != nil {
		server.log.Errorf("Get event from request body: %v", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if event.Title == "" || event.StartTime.IsZero() || event.UserID == 0 {
		errMsg := "Empty fields event: title, start time, user id"
		server.log.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := server.app.CreateEvent(req.Context(), event); err != nil {
		server.log.Errorf("Saving event to storage: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func handleGetEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	id, err := getID(req)
	if err != nil {
		server.log.Errorf("Get event id from request body: %v", err)
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	event, err := server.app.GetEvent(req.Context(), id)
	if err != nil {
		server.log.Errorf("Get event from storage: %v", err)
		if errors.Is(err, storage.ErrNoEvent) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		server.log.Errorf("Marshaling event to json: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if _, err := res.Write(jsonEvent); err != nil {
		server.log.Errorf("Write event to response: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUpdateEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	// Update only required fields
	event, err := unmarshalEvent(req)
	if err != nil {
		server.log.Errorf("Get event from request body: %v", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if event.ID == 0 {
		errMsg := "Event id not set"
		server.log.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := server.app.UpdateEvent(req.Context(), event); err != nil {
		server.log.Errorf("Update event in storage: %v", err)
		if errors.Is(err, storage.ErrNoEvent) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func handleDeleteEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	id, err := getID(req)
	if err != nil {
		server.log.Errorf("Get event id from request body: %v", err)
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := server.app.DeleteEvent(req.Context(), id); err != nil {
		server.log.Errorf("Delete event in storage: %v", err)
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
