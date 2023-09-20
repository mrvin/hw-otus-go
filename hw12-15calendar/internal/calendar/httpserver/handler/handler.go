package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/app"
)

var ErrIDEmpty = errors.New("id is empty")

type Handler struct {
	app *app.App
}

func New(a *app.App) *Handler {
	return &Handler{
		app: a,
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
