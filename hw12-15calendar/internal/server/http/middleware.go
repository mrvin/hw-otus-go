package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

var ctx = context.Background()

func handleCreateUser(res http.ResponseWriter, req *http.Request, server *Server) {
	user, err := unmarshalUser(req)
	if err != nil {
		server.logg.Print(err)
		return
	}
	if user.Name == "" || user.Email == "" {
		res.WriteHeader(http.StatusBadRequest)
	}

	if err := server.stor.CreateUser(ctx, user); err != nil {
		server.logg.Print(err)
		return
	}

	res.WriteHeader(200)
}

func handleGetUser(res http.ResponseWriter, req *http.Request, server *Server) {
	id, err := getID(req)
	if err != nil {
		server.logg.Print(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := server.stor.GetUser(ctx, id)
	if err != nil {
		server.logg.Print(err)
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		server.logg.Printf("can't marshaling json: %v", err)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonUser)
}

func handleUpdateUser(res http.ResponseWriter, req *http.Request, server *Server) {
	user, err := unmarshalUser(req)
	if err != nil {
		server.logg.Print(err)
		return
	}
	if user.ID == 0 {
		res.WriteHeader(http.StatusBadRequest)
	}

	if err := server.stor.UpdateUser(ctx, user); err != nil {
		server.logg.Print(err)
	}
}

func handleDeleteUser(res http.ResponseWriter, req *http.Request, server *Server) {
	id, _ := getID(req)

	if err := server.stor.DeleteUser(ctx, id); err != nil {
		server.logg.Print(err)
	}
}

func getID(req *http.Request) (id int, err error) {
	idStr := req.URL.Query().Get("id")
	id, err = strconv.Atoi(idStr)
	if err != nil {
		return
	}

	return
}

func unmarshalUser(req *http.Request) (*storage.User, error) {
	var user storage.User

	lenReq := req.ContentLength
	body := make([]byte, lenReq)
	req.Body.Read(body)

	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
