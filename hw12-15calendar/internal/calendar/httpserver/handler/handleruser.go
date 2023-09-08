package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

// TODO:add return id created user.
func (h *Handler) CreateUser(res http.ResponseWriter, req *http.Request) {
	user, err := unmarshalUser(req)
	if err != nil {
		h.log.Errorf("Get user from request body: %v", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Name == "" || user.Email == "" {
		errMsg := "Empty fields user: name, email"
		h.log.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := h.app.CreateUser(req.Context(), user); err != nil {
		h.log.Errorf("Saving user to storage: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

//nolint:dupl
func (h *Handler) GetUser(res http.ResponseWriter, req *http.Request) {
	id, err := getID(req)
	if err != nil {
		h.log.Errorf("Get user id from request body: %v", err)
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	user, err := h.app.GetUser(req.Context(), id)
	if err != nil {
		h.log.Errorf("Get user from storage: %v", err)
		if errors.Is(err, storage.ErrNoUser) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		h.log.Errorf("Marshaling user to json: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if _, err := res.Write(jsonUser); err != nil {
		h.log.Errorf("Write user to response: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateUser(res http.ResponseWriter, req *http.Request) {
	// Update only required fields
	user, err := unmarshalUser(req)
	if err != nil {
		h.log.Errorf("Get user from request body: %v", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if user.ID == 0 {
		errMsg := "User id not set"
		h.log.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := h.app.UpdateUser(req.Context(), user); err != nil {
		h.log.Errorf("Update user in storage: %v", err)
		if errors.Is(err, storage.ErrNoUser) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func (h *Handler) DeleteUser(res http.ResponseWriter, req *http.Request) {
	id, err := getID(req)
	if err != nil {
		h.log.Errorf("Get user id from request body: %v", err)
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := h.app.DeleteUser(req.Context(), id); err != nil {
		h.log.Errorf("Delete user in storage: %v", err)
		if errors.Is(err, storage.ErrNoUser) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func unmarshalUser(req *http.Request) (*storage.User, error) {
	var user storage.User

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read body req: %w", err)
	}

	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("unmarshal body req: %w", err)
	}

	return &user, nil
}
