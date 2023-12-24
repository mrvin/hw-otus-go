package handleruser

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar/server/http/handlers"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	httpresponse "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/response"
)

//go:generate go run github.com/vektra/mockery/v2@v2.38.0 --name=UserDeleter
type UserDeleter interface {
	DeleteUser(ctx context.Context, name string) error
}

func New(deleter UserDeleter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userName := handler.GetUserName(req.Context())
		if userName == "" {
			err := fmt.Errorf("DeleteUser: user name is empty")
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := deleter.DeleteUser(req.Context(), userName); err != nil {
			err := fmt.Errorf("DeleteUser: delete user from storage: %w", err)
			slog.Error(err.Error())
			if errors.Is(err, storage.ErrNoUser) {
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			} else {
				httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Write json response
		httpresponse.WriteOK(res)

		slog.Info("User deletion was successful")
	}
}
