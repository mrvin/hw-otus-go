package handler

import (
	"github.com/mrvin/hw-otus-go/hw12-15calendar/services/calendar/app"
	"go.uber.org/zap"
)

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
