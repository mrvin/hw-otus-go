package handlerevent

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	httpresponse "github.com/mrvin/hw-otus-go/hw12-15calendar/pkg/http/response"
)

func (h *Handler) DeleteEvent(res http.ResponseWriter, req *http.Request) {
	id, err := getID(req)
	if err != nil {
		err := fmt.Errorf("DeleteEvent: get event id from request url query: %w", err)
		slog.Error(err.Error())
		if errors.Is(err, ErrIDEmpty) {
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		} else {
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := h.eventService.DeleteEvent(req.Context(), id); err != nil {
		err := fmt.Errorf("DeleteEvent: delete event from storage: %w", err)
		slog.Error(err.Error())
		if errors.Is(err, storage.ErrNoEvent) {
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
		} else {
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Write json response
	httpresponse.WriteOK(res)

	slog.Info("Event deletion was successful")
}
