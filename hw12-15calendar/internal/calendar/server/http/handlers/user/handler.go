package handleruser

import (
	authservice "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/service/auth"
)

type Handler struct {
	authService *authservice.AuthService
}

func New(auth *authservice.AuthService) *Handler {
	return &Handler{
		authService: auth,
	}
}
