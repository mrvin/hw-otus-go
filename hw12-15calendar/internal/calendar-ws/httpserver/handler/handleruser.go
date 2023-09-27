package handler

import (
	"log/slog"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DisplayListUsers(res http.ResponseWriter, req *http.Request) {
	users, err := h.grpcclient.GetAllUsers(req.Context(), &emptypb.Empty{})
	if err != nil {
		slog.Error("gRPC get all users: GetAllUsers: " + err.Error())
		h.ErrMsg(res)
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
		slog.Error("Execute display list users template: " + err.Error())
		h.ErrMsg(res)
		return
	}
}

func (h *Handler) DisplayUser(res http.ResponseWriter, req *http.Request) {
	idUser, err := getID(req)
	if err != nil {
		slog.Error("Get id user: " + err.Error())
		h.ErrMsg(res)
		return
	}
	reqUser := &calendarapi.UserRequest{Id: int64(idUser)}
	user, err := h.grpcclient.GetUser(req.Context(), reqUser)
	if err != nil {
		slog.Error("gRPC get user: " + err.Error())
		h.ErrMsg(res)
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
		slog.Error("Execute display user template: " + err.Error())
		h.ErrMsg(res)
		return
	}

}

func (h *Handler) DisplayFormUser(res http.ResponseWriter, req *http.Request) {
	data := resp{Title: "Create user"}
	if err := h.templates.Execute("form-user.html", res, data); err != nil {
		slog.Error("Execute display form user template: " + err.Error())
		h.ErrMsg(res)
		return
	}
}

func (h *Handler) CreateUser(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		slog.Error("Parse form сreate user: " + err.Error())
		h.ErrMsg(res)
		return
	}

	name := req.FormValue("name")
	email := req.FormValue("email")

	user := &calendarapi.User{Name: name, Email: email}
	_, err := h.grpcclient.CreateUser(req.Context(), user)
	if err != nil {
		slog.Error("gRPC create user: " + err.Error())
		h.ErrMsg(res)
		return
	}
	h.DisplayListUsers(res, req)
}

func (h *Handler) DeleteUser(res http.ResponseWriter, req *http.Request) {
	idUser, err := getID(req)
	if err != nil {
		slog.Error("Get id user: " + err.Error())
		h.ErrMsg(res)
		return
	}
	reqUser := &calendarapi.UserRequest{Id: int64(idUser)}
	if _, err := h.grpcclient.DeleteUser(req.Context(), reqUser); err != nil {
		slog.Error("gRPC delete user: ", err.Error())
		h.ErrMsg(res)
		return
	}

	text := resp{Title: "Delete user", Body: struct{ Text string }{"User deleted successfully"}}
	if err := h.templates.Execute("text.html", res, text); err != nil {
		slog.Error("Execute delete user template: " + err.Error())
		h.ErrMsg(res)
		return
	}
}
