package httpserver

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ErrIDEmpty = errors.New("id is empty")

type resp struct {
	Title string
	Body  struct {
		Text string
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

func errMsg(res http.ResponseWriter, templates *templateLoader) {
	text := resp{Title: "Error", Body: struct{ Text string }{"Sorry something went wrong."}}
	if err := executeTemp("text.html", res, text, templates); err != nil {
		log.Printf("errMsg: %v", err)
		return
	}
}

func executeTemp(nameTemp string, res http.ResponseWriter, data any, templates *templateLoader) error {
	temp, ok := templates.templates[nameTemp]
	if !ok {
		return fmt.Errorf("not found template '%s'", nameTemp)
	}

	if err := temp.Execute(res, data); err != nil {
		return fmt.Errorf("execute %w", err)
	}

	return nil
}

func displayListUsers(res http.ResponseWriter, req *http.Request, server *Server) {
	users, err := server.grpcclient.GetAllUsers(req.Context(), &emptypb.Empty{})
	if err != nil {
		log.Printf("displayListUsers: GetAllUsers: %v", err)
		errMsg(res, server.templates)
		return
	}

	dataUsers := struct {
		Title string
		Body  struct {
			Users []*calendarapi.User
		}
	}{
		Title: "List users",
		Body: struct {
			Users []*calendarapi.User
		}{users.Users},
	}

	if err := executeTemp("list-users.html", res, dataUsers, server.templates); err != nil {
		log.Printf("displayListUsers: %v", err)
		errMsg(res, server.templates)
		return
	}
}

func displayUser(res http.ResponseWriter, req *http.Request, server *Server) {
	idUser, err := getID(req)
	if err != nil {
		log.Printf("Get id user: %v", err)
		errMsg(res, server.templates)
		return
	}
	reqUser := &calendarapi.UserRequest{Id: int64(idUser)}
	user, err := server.grpcclient.GetUser(req.Context(), reqUser)
	if err != nil {
		log.Printf("displayUser: %v", err)
		errMsg(res, server.templates)
		return
	}

	dataUser := struct {
		Title string
		Body  struct {
			User *calendarapi.User
		}
	}{
		Title: "User",
		Body: struct {
			User *calendarapi.User
		}{user},
	}

	if err := executeTemp("user.html", res, dataUser, server.templates); err != nil {
		log.Printf("displayUser: %v", err)
		errMsg(res, server.templates)
		return
	}

}

func displayFormUser(res http.ResponseWriter, req *http.Request, server *Server) {
	data := resp{Title: "Create user"}
	if err := executeTemp("form-user.html", res, data, server.templates); err != nil {
		log.Printf("displayFormUser: %v", err)
		errMsg(res, server.templates)
		return
	}
}

func createUser(res http.ResponseWriter, req *http.Request, server *Server) {
	if err := req.ParseForm(); err != nil {
		log.Printf("сreateUser: %v", err)
		errMsg(res, server.templates)
		return
	}

	name := req.FormValue("name")
	email := req.FormValue("email")

	user := &calendarapi.User{Name: name, Email: email}
	_, err := server.grpcclient.CreateUser(req.Context(), user)
	if err != nil {
		log.Printf("сreateUser: %v", err)
		errMsg(res, server.templates)
		return
	}

	text := resp{Title: "Create user", Body: struct{ Text string }{"User created successfully"}}
	if err := executeTemp("text.html", res, text, server.templates); err != nil {
		log.Printf("сreateUser: %v", err)
		errMsg(res, server.templates)
		return
	}
}

func deleteUser(res http.ResponseWriter, req *http.Request, server *Server) {
	idUser, err := getID(req)
	if err != nil {
		log.Printf("Get id user: %v", err)
		errMsg(res, server.templates)
		return
	}
	reqUser := &calendarapi.UserRequest{Id: int64(idUser)}
	if _, err := server.grpcclient.DeleteUser(req.Context(), reqUser); err != nil {
		log.Printf("Delete user: %v", err)
		errMsg(res, server.templates)
		return
	}

	text := resp{Title: "Delete user", Body: struct{ Text string }{"User deleted successfully"}}
	if err := executeTemp("text.html", res, text, server.templates); err != nil {
		log.Printf("Delete user: %v", err)
		errMsg(res, server.templates)
		return
	}
}
