package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar-ws/httpserver/templateloader"
)

var ErrIDEmpty = errors.New("id is empty")

type Handler struct {
	templates *templateloader.TemplateLoader
	grpcclient calendarapi.EventServiceClient
	log       *zap.SugaredLogger
}

func New(templates *templateloader.TemplateLoader, grpcclient calendarapi.EventServiceClient, log *zap.SugaredLogger) *Handler {
	templates.Load("templates")
	return &Handler{
		templates: templates,
		grpcclient: grpcclient,
		log: log,
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

func errMsg(res http.ResponseWriter, templates *templateloader.TemplateLoader) {
	text := resp{Title: "Error", Body: struct{ Text string }{"Sorry something went wrong."}}
	if err := templates.Execute("text.html", res, text); err != nil {
		log.Printf("errMsg: %v", err)
		return
	}
}
