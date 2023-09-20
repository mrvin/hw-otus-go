package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			UserID int
		}
	}{
		Title: "Create event",
		Body: struct {
			UserID int
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
	idUser, err := strconv.Atoi(idStr)
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

	starTimePB := timestamppb.New(startTimeGO)
	stopTimePB := timestamppb.New(stopTimeGO)
	event := &calendarapi.Event{Title: title, Description: description, StartTime: starTimePB, StopTime: stopTimePB, UserID: int64(idUser)}
	_, err = h.grpcclient.CreateEvent(req.Context(), event)
	if err != nil {
		slog.Error("gRPC сreate event: " + err.Error())
		h.ErrMsg(res)
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
}

func (h *Handler) DisplayEvent(res http.ResponseWriter, req *http.Request) {
	idEvent, err := getID(req)
	if err != nil {
		slog.Error("Get id event: " + err.Error())
		h.ErrMsg(res)
		return
	}
	reqEvent := &calendarapi.EventRequest{Id: int64(idEvent)}
	event, err := h.grpcclient.GetEvent(req.Context(), reqEvent)
	if err != nil {
		slog.Error("gRPC get event: " + err.Error())
		h.ErrMsg(res)
		return
	}

	storageEvent := storage.Event{ID: int(event.Id), Title: event.Title,
		Description: event.Description, StartTime: event.StartTime.AsTime(), StopTime: event.StopTime.AsTime()}

	dataEvent := struct {
		Title string
		Body  struct {
			Event storage.Event
		}
	}{
		Title: "Event",
		Body: struct {
			Event storage.Event
		}{storageEvent},
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

	reqUser := &calendarapi.GetEventsForUserRequest{
		User: &calendarapi.UserRequest{Id: int64(idUser)},
		DaysAhead: &calendarapi.DaysAheadRequest{
			Days: int32(days),
			Date: timestamppb.New(time.Now())}}

	events, err := h.grpcclient.GetEventsForUser(req.Context(), reqUser)
	if err != nil {
		slog.Error("gRPC list events for user: " + err.Error())
		h.ErrMsg(res)
		return
	}

	dataEvent := struct {
		Title string
		Body  struct {
			UserID int
			Events []*calendarapi.Event
		}
	}{
		Title: "List events",
		Body: struct {
			UserID int
			Events []*calendarapi.Event
		}{idUser, events.Events},
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
	reqEvent := &calendarapi.EventRequest{Id: int64(idEvent)}
	if _, err := h.grpcclient.DeleteEvent(req.Context(), reqEvent); err != nil {
		slog.Error("gRPC delete event: " + err.Error())
		h.ErrMsg(res)
		return
	}

	text := resp{Title: "Delete event", Body: struct{ Text string }{"Event deleted successfully"}}
	if err := h.templates.Execute("text.html", res, text); err != nil {
		slog.Error("Execute delete event template: " + err.Error())
		h.ErrMsg(res)
		return
	}
}
