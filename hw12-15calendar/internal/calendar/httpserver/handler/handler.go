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

type responseID struct {
	ID int64 `json:"id"`
}

func getID(req *http.Request) (int64, error) {
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		return 0, ErrIDEmpty
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("convert id: %w", err)
	}

	return id, nil
}
