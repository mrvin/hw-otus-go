package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/cmd/calendar/app"
	"go.uber.org/zap"
)

var ErrIDEmpty = errors.New("id is empty")

type Handler struct {
	app *app.App
	log *zap.SugaredLogger
}

func New(a *app.App, log *zap.SugaredLogger) *Handler {
	return &Handler{
		app: a,
		log: log,
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
