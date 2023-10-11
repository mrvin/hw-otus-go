package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

//TODO: http.Error(res, err.Error(), http.StatusBadRequest)

func (h *Handler) DisplayFormEvent(res http.ResponseWriter, req *http.Request) {
	idUser, err := getID(req)
	if err != nil {
		slog.Error("Get id user: " + err.Error())
		h.ErrMsg(res)
		return
	}

	data := struct {
		Title string
		Body  struct {
			UserID int64
		}
	}{
		Title: "Create event",
		Body: struct {
			UserID int64
		}{idUser},
	}
	if err := h.templates.Execute("form-event.html", res, data); err != nil {
		slog.Error("Execute display form event template: " + err.Error())
		h.ErrMsg(res)
		return
	}
}

func (h *Handler) CreateEvent(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		slog.Error("Parse create event form: " + err.Error())
		h.ErrMsg(res)
		return
	}

	idStr := req.FormValue("id")
	title := req.FormValue("title")
	description := req.FormValue("description")
	startTime := req.FormValue("startTime")
	stopTime := req.FormValue("stopTime")
	timeZone := req.FormValue("timezone")

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		slog.Error("Get location: " + err.Error())
		h.ErrMsg(res)
		return
	}
	idUser, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("Convert id: " + err.Error())
		h.ErrMsg(res)
		return
	}
	timeNow := time.Now().In(loc)
	tZone := timeNow.Format("Z07:00")
	startTimeGO, err := time.ParseInLocation(time.RFC3339, startTime+":00"+tZone, loc)
	if err != nil {
		slog.Error("Parse start time location: " + err.Error())
		h.ErrMsg(res)
		return
	}

	stopTimeGO, err := time.ParseInLocation(time.RFC3339, stopTime+":00"+tZone, loc)
	if err != nil {
		slog.Error("Parse stop time location: " + err.Error())
		h.ErrMsg(res)
		return
	}

	if _, err := h.client.CreateEvent(req.Context(), title, description, startTimeGO, stopTimeGO, idUser); err != nil {
		slog.Error("Create event: " + err.Error())
		return
	}

	reqListEvents, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/list-events?id=%d&days=%d", idUser, 0), nil)
	if err != nil {
		slog.Error("Request list events: " + err.Error())
		h.ErrMsg(res)
		return
	}
	// FIXIT path: POST /create-event
	h.DisplayListEventsForUser(res, reqListEvents)
	//http.Redirect(res, reqListEvents, reqListEvents.URL.RequestURI(), http.StatusTemporaryRedirect)
}

func (h *Handler) DisplayEvent(res http.ResponseWriter, req *http.Request) {
	idEvent, err := getID(req)
	if err != nil {
		slog.Error("Get id event: " + err.Error())
		h.ErrMsg(res)
		return
	}

	event, err := h.client.GetEvent(req.Context(), idEvent)
	if err != nil {
		slog.Error("Get event: " + err.Error())
		return
	}

	dataEvent := struct {
		Title string
		Body  struct {
			Event storage.Event
		}
	}{
		Title: "Event",
		Body: struct {
			Event storage.Event
		}{*event},
	}

	if err := h.templates.Execute("event.html", res, dataEvent); err != nil {
		slog.Error("Execute event template: " + err.Error())
		h.ErrMsg(res)
		return
	}
}

func (h *Handler) DisplayListEventsForUser(res http.ResponseWriter, req *http.Request) {
	idUser, err := getID(req)
	if err != nil {
		slog.Error("Get id user: " + err.Error())
		h.ErrMsg(res)
		return
	}
	days := 0
	daysStr := req.URL.Query().Get("days")
	if daysStr != "" {
		days, err = strconv.Atoi(daysStr)
		if err != nil {
			slog.Error("Convert days: " + err.Error())
			return
		}
	}

	events, err := h.client.ListEventsForUser(req.Context(), idUser, days)

	dataEvent := struct {
		Title string
		Body  struct {
			UserID int64
			Events []storage.Event
		}
	}{
		Title: "List events",
		Body: struct {
			UserID int64
			Events []storage.Event
		}{idUser, events},
	}

	if err := h.templates.Execute("list-events.html", res, dataEvent); err != nil {
		slog.Error("Execute display list events template: " + err.Error())
		h.ErrMsg(res)
		return
	}
}

func (h *Handler) DeleteEvent(res http.ResponseWriter, req *http.Request) {
	idEvent, err := getID(req)
	if err != nil {
		slog.Error("Get id event: " + err.Error())
		h.ErrMsg(res)
		return
	}

	if err := h.client.DeleteEvent(req.Context(), idEvent); err != nil {
		slog.Error("Delete event: " + err.Error())
		return
	}

	text := resp{Title: "Delete event", Body: struct{ Text string }{"Event deleted successfully"}}
	if err := h.templates.Execute("text.html", res, text); err != nil {
		slog.Error("Execute delete event template: " + err.Error())
		h.ErrMsg(res)
		return
	}
}
