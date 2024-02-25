package handlers

import (
	"log/slog"
	"net/http"
)

func (h *Handler) DisplayFormRegistration(res http.ResponseWriter, req *http.Request) {
	data := resp{Title: "Create user"}
	if err := h.templates.Execute("form-user.html", res, data); err != nil {
		slog.Error("Execute display form user template: " + err.Error())
		h.ErrMsg(res)
		return
	}
}

func (h *Handler) Registration(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		slog.Error("Parse form сreate user: " + err.Error())
		h.ErrMsg(res)
		return
	}

	name := req.FormValue("name")
	password := req.FormValue("password")
	email := req.FormValue("email")

	if _, err := h.client.Registration(req.Context(), name, password, email); err != nil {
		slog.Error("Create user: " + err.Error())
		return
	}

	urlLogin := "/form-login"
	var err error
	req, err = http.NewRequest(http.MethodGet, urlLogin, nil)
	if err != nil {
		slog.Error("Create new request: " + err.Error())
		return
	}

	http.Redirect(res, req, urlLogin, http.StatusSeeOther)
}
