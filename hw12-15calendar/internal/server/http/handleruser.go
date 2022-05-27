package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

var ctx = context.Background()

var ErrIDEmpty = errors.New("id is empty")

func handleCreateUser(res http.ResponseWriter, req *http.Request, server *Server) {
	user, err := unmarshalUser(req)
	if err != nil {
		log.Printf("Create user: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if user.Name == "" || user.Email == "" {
		log.Printf("Create user: empty name or email: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := server.stor.CreateUser(ctx, user); err != nil {
		log.Printf("Create user: saving user info: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func handleGetUser(res http.ResponseWriter, req *http.Request, server *Server) {
	// Return an error in JSON
	id, err := getID(req)
	if err != nil {
		log.Printf("Get user: get id: %v", err)
		if err == ErrIDEmpty { //nolint:errorlint
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	user, err := server.stor.GetUser(ctx, id)
	if err != nil {
		log.Printf("Get user: get from storage: %v", err)
		if errors.Is(err, storage.ErrNoUser) {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Printf("Get user: marshaling json: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	if _, err := res.Write(jsonUser); err != nil {
		log.Printf("Get user: write res: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleUpdateUser(res http.ResponseWriter, req *http.Request, server *Server) {
	// Update only required fields
	user, err := unmarshalUser(req)
	if err != nil {
		log.Printf("Update user: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if user.ID == 0 {
		log.Print("Update user: user id not set")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := server.stor.UpdateUser(ctx, user); err != nil {
		log.Printf("Update user: update in storage: %v", err)
		if errors.Is(err, storage.ErrNoUser) {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}

func handleDeleteUser(res http.ResponseWriter, req *http.Request, server *Server) {
	id, err := getID(req)
	if err != nil {
		log.Printf("Delete user: get id: %v", err)
		if err == ErrIDEmpty { //nolint:errorlint
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if err := server.stor.DeleteUser(ctx, id); err != nil {
		log.Printf("Delete user: delete in storage: %v", err)
		if errors.Is(err, storage.ErrNoUser) {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
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

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("read body req: %w", err)
	}

	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("unmarshal body req: %w", err)
	}

	return &user, nil
}
