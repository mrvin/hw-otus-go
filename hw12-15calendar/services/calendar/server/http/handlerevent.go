package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func handleCreateEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	event, err := unmarshalEvent(req)
	if err != nil {
		log.Printf("Create event: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if event.Title == "" || event.StartTime.IsZero() || event.UserID == 0 {
		log.Print("Create event: empty title or start time or user id")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := server.stor.CreateEvent(ctx, event); err != nil {
		log.Printf("Create event: saving event info: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func handleGetEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	// Return an error in JSON
	id, err := getID(req)
	if err != nil {
		log.Printf("Get event: get event id: %v", err)
		if errors.Is(err, ErrIDEmpty) {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	event, err := server.stor.GetEvent(ctx, id)
	if err != nil {
		log.Printf("Get event: get from storage: %v", err)
		if errors.Is(err, storage.ErrNoEvent) {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		log.Printf("Get event: marshaling json: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if _, err := res.Write(jsonEvent); err != nil {
		log.Printf("Get event: write res: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleUpdateEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	// Update only required fields
	event, err := unmarshalEvent(req)
	if err != nil {
		log.Printf("Update event: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if event.ID == 0 {
		log.Print("Update event: event id not set")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := server.stor.UpdateEvent(ctx, event); err != nil {
		log.Printf("Update event: update in storage: %v", err)
		if errors.Is(err, storage.ErrNoEvent) {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}

func handleDeleteEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	id, err := getID(req)
	if err != nil {
		log.Printf("Delete event: get id: %v", err)
		if errors.Is(err, ErrIDEmpty) {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if err := server.stor.DeleteEvent(ctx, id); err != nil {
		log.Printf("Delete event: delete in storage: %v", err)
		if errors.Is(err, storage.ErrNoEvent) {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}

func unmarshalEvent(req *http.Request) (*storage.Event, error) {
	var event storage.Event

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("read body req: %w", err)
	}

	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("unmarshal body req: %w", err)
	}

	return &event, nil
}
