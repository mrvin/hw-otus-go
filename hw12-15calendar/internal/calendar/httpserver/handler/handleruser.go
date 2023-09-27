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

// TODO:add return id created user.
func (h *Handler) CreateUser(res http.ResponseWriter, req *http.Request) {
	user, err := unmarshalUser(req)
	if err != nil {
		slog.Error("Get user from request body: " + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Name == "" || user.Email == "" {
		errMsg := "Empty fields user: name, email"
		slog.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := h.app.CreateUser(req.Context(), user); err != nil {
		slog.Error("Saving user to storage: " + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

//nolint:dupl
func (h *Handler) GetUser(res http.ResponseWriter, req *http.Request) {
	id, err := getID(req)
	if err != nil {
		slog.Error("Get user id from request body: " + err.Error())
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	user, err := h.app.GetUser(req.Context(), id)
	if err != nil {
		slog.Error("Get user from storage: " + err.Error())
		if errors.Is(err, storage.ErrNoUser) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		slog.Error("Marshaling user to json: " + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if _, err := res.Write(jsonUser); err != nil {
		slog.Error("Write user to response: " + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateUser(res http.ResponseWriter, req *http.Request) {
	// Update only required fields
	user, err := unmarshalUser(req)
	if err != nil {
		slog.Error("Get user from request body: " + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if user.ID == 0 {
		errMsg := "User id not set"
		slog.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := h.app.UpdateUser(req.Context(), user); err != nil {
		slog.Error("Update user in storage: " + err.Error())
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
		slog.Error("Get user id from request body: " + err.Error())
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := h.app.DeleteUser(req.Context(), id); err != nil {
		slog.Error("Delete user in storage: " + err.Error())
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
