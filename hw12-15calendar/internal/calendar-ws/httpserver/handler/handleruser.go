package handler

import (
	"log/slog"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func (h *Handler) DisplayListUsers(res http.ResponseWriter, req *http.Request) {
	users, err := h.client.ListUsers(req.Context())
	if err != nil {
		slog.Error("List users: " + err.Error())
		return
	}

	dataUsers := struct {
		Title string
		Body  struct {
			Users []storage.User
		}
	}{
		Title: "List users",
		Body: struct {
			Users []storage.User
		}{users},
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
	user, err := h.client.GetUser(req.Context(), idUser)
	if err != nil {
		slog.Error("Get user: " + err.Error())
		return
	}

	dataUser := struct {
		Title string
		Body  struct {
			User *storage.User
		}
	}{
		Title: "User",
		Body: struct {
			User *storage.User
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
	password := req.FormValue("password")
	email := req.FormValue("email")

	if _, err := h.client.CreateUser(req.Context(), name, password, email); err != nil {
		slog.Error("Create user: " + err.Error())
		return
	}

	// FIXIT
	h.DisplayListUsers(res, req)
}

func (h *Handler) DeleteUser(res http.ResponseWriter, req *http.Request) {
	userName := req.URL.Query().Get("name")
	if userName == "" {
		slog.Error("Cant get user name")
		h.ErrMsg(res)
		return
	}
	if err := h.client.DeleteUser(req.Context(), userName); err != nil {
		slog.Error("Delete user: " + err.Error())
		return
	}

	text := resp{Title: "Delete user", Body: struct{ Text string }{"User deleted successfully"}}
	if err := h.templates.Execute("text.html", res, text); err != nil {
		slog.Error("Execute delete user template: " + err.Error())
		h.ErrMsg(res)
		return
	}
}
