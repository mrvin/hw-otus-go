package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

var ErrIDEmpty = errors.New("id is empty")

// TODO:add return id created user
func handleCreateUser(res http.ResponseWriter, req *http.Request, server *Server) {
	user, err := unmarshalUser(req)
	if err != nil {
		server.log.Errorf("Get user from request body: %v", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Name == "" || user.Email == "" {
		errMsg := "Empty fields user: name, email"
		server.log.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := server.app.CreateUser(req.Context(), user); err != nil {
		server.log.Errorf("Saving user to storage: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func handleGetUser(res http.ResponseWriter, req *http.Request, server *Server) {
	id, err := getID(req)
	if err != nil {
		server.log.Errorf("Get user id from request body: %v", err)
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	user, err := server.app.GetUser(req.Context(), id)
	if err != nil {
		server.log.Errorf("Get user from storage: %v", err)
		if errors.Is(err, storage.ErrNoUser) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		server.log.Errorf("Marshaling user to json: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if _, err := res.Write(jsonUser); err != nil {
		server.log.Errorf("Write user to response: %v", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUpdateUser(res http.ResponseWriter, req *http.Request, server *Server) {
	// Update only required fields
	user, err := unmarshalUser(req)
	if err != nil {
		server.log.Errorf("Get user from request body: %v", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if user.ID == 0 {
		errMsg := "User id not set"
		server.log.Error(errMsg)
		http.Error(res, errMsg, http.StatusBadRequest)
		return
	}

	if err := server.app.UpdateUser(req.Context(), user); err != nil {
		server.log.Errorf("Update user in storage: %v", err)
		if errors.Is(err, storage.ErrNoUser) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func handleDeleteUser(res http.ResponseWriter, req *http.Request, server *Server) {
	id, err := getID(req)
	if err != nil {
		server.log.Errorf("Get user id from request body: %v", err)
		if errors.Is(err, ErrIDEmpty) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := server.app.DeleteUser(req.Context(), id); err != nil {
		server.log.Errorf("Delete user in storage: %v", err)
		if errors.Is(err, storage.ErrNoUser) {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func getID(req *http.Request) (int, error) {
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		return 0, ErrIDEmpty
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("convert id: %w", err)
	}

	return id, nil
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
