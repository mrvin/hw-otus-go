package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-ws/httpserver/templateloader"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-api"
)

var ErrIDEmpty = errors.New("id is empty")

type Handler struct {
	templates  *templateloader.TemplateLoader
	grpcclient calendarapi.EventServiceClient
}

func New(grpcclient calendarapi.EventServiceClient) *Handler {
	templates := templateloader.New()
	templates.Load("templates")
	return &Handler{
		templates:  templates,
		grpcclient: grpcclient,
	}
}

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

func (h *Handler) ErrMsg(res http.ResponseWriter) {
	text := resp{Title: "Error", Body: struct{ Text string }{"Sorry something went wrong."}}
	if err := h.templates.Execute("text.html", res, text); err != nil {
		slog.Error("Execut error tetemplate: " + err.Error())
		return
	}
}
