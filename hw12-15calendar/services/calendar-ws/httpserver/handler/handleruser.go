package handler

import (
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DisplayListUsers(res http.ResponseWriter, req *http.Request) {
	users, err := h.grpcclient.GetAllUsers(req.Context(), &emptypb.Empty{})
	if err != nil {
		h.log.Errorf("displayListUsers: GetAllUsers: %v", err)
		errMsg(res, h.templates)
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

	if err := h.templates.Execute("list-users.html", res, dataUsers); err != nil {
		h.log.Errorf("displayListUsers: %v", err)
		errMsg(res, h.templates)
		return
	}
}

func (h *Handler) DisplayUser(res http.ResponseWriter, req *http.Request) {
	idUser, err := getID(req)
	if err != nil {
		h.log.Errorf("Get id user: %v", err)
		errMsg(res, h.templates)
		return
	}
	reqUser := &calendarapi.UserRequest{Id: int64(idUser)}
	user, err := h.grpcclient.GetUser(req.Context(), reqUser)
	if err != nil {
		h.log.Errorf("displayUser: %v", err)
		errMsg(res, h.templates)
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

	if err := h.templates.Execute("user.html", res, dataUser); err != nil {
		h.log.Errorf("displayUser: %v", err)
		errMsg(res, h.templates)
		return
	}

}

func (h *Handler) DisplayFormUser(res http.ResponseWriter, req *http.Request) {
	data := resp{Title: "Create user"}
	if err := h.templates.Execute("form-user.html", res, data); err != nil {
		h.log.Errorf("displayFormUser: %v", err)
		errMsg(res, h.templates)
		return
	}
}

func (h *Handler) CreateUser(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		h.log.Errorf("сreateUser: %v", err)
		errMsg(res, h.templates)
		return
	}

	name := req.FormValue("name")
	email := req.FormValue("email")

	user := &calendarapi.User{Name: name, Email: email}
	_, err := h.grpcclient.CreateUser(req.Context(), user)
	if err != nil {
		h.log.Errorf("сreateUser: %v", err)
		errMsg(res, h.templates)
		return
	}
	h.DisplayListUsers(res, req)
}

func (h *Handler) DeleteUser(res http.ResponseWriter, req *http.Request) {
	idUser, err := getID(req)
	if err != nil {
		h.log.Errorf("Get id user: %v", err)
		errMsg(res, h.templates)
		return
	}
	reqUser := &calendarapi.UserRequest{Id: int64(idUser)}
	if _, err := h.grpcclient.DeleteUser(req.Context(), reqUser); err != nil {
		h.log.Errorf("Delete user: %v", err)
		errMsg(res, h.templates)
		return
	}

	text := resp{Title: "Delete user", Body: struct{ Text string }{"User deleted successfully"}}
	if err := h.templates.Execute("text.html", res, text); err != nil {
		h.log.Errorf("Delete user: %v", err)
		errMsg(res, h.templates)
		return
	}
}
