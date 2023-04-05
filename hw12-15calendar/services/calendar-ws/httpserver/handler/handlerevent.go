package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) DisplayFormEvent(res http.ResponseWriter, req *http.Request) {
	idUser, err := getID(req)
	if err != nil {
		h.log.Errorf("Get id user: %v", err)
		errMsg(res, h.templates)
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
		h.log.Errorf("displayFormEvent: %v", err)
		errMsg(res, h.templates)
		return
	}
}

func (h *Handler) CreateEvent(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		h.log.Errorf("createEvent: %v", err)
		errMsg(res, h.templates)
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
		h.log.Errorf("get location: %v", err)
		errMsg(res, h.templates)
		return
	}
	idUser, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Errorf("convert id: %v", err)
		errMsg(res, h.templates)
		return
	}
	timeNow := time.Now().In(loc)
	tZone := timeNow.Format("Z07:00")
	startTimeGO, err := time.ParseInLocation(time.RFC3339, startTime+":00"+tZone, loc)
	if err != nil {
		h.log.Errorf("сreateEvent: parse startTime: %v", err)
		errMsg(res, h.templates)
		return
	}

	stopTimeGO, err := time.ParseInLocation(time.RFC3339, stopTime+":00"+tZone, loc)
	if err != nil {
		h.log.Errorf("сreateEvent: parse stopTime: %v", err)
		errMsg(res, h.templates)
		return
	}

	starTimePB := timestamppb.New(startTimeGO)
	stopTimePB := timestamppb.New(stopTimeGO)
	event := &calendarapi.Event{Title: title, Description: description, StartTime: starTimePB, StopTime: stopTimePB, UserID: int64(idUser)}
	_, err = h.grpcclient.CreateEvent(req.Context(), event)
	if err != nil {
		h.log.Errorf("сreateEvent: %v", err)
		errMsg(res, h.templates)
		return
	}

	text := resp{Title: "Create event", Body: struct{ Text string }{"Event created successfully"}}
	if err := h.templates.Execute("text.html", res, text); err != nil {
		h.log.Errorf("createEvent: %v", err)
		errMsg(res, h.templates)
		return
	}
}

func (h *Handler) DisplayEvent(res http.ResponseWriter, req *http.Request) {
	idEvent, err := getID(req)
	if err != nil {
		h.log.Errorf("Get id event: %v", err)
		errMsg(res, h.templates)
		return
	}
	reqEvent := &calendarapi.EventRequest{Id: int64(idEvent)}
	event, err := h.grpcclient.GetEvent(req.Context(), reqEvent)
	if err != nil {
		h.log.Errorf("displayEvent: %v", err)
		errMsg(res, h.templates)
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
		h.log.Errorf("displayEvent: %v", err)
		errMsg(res, h.templates)
		return
	}
}

func (h *Handler) DisplayListEventsForUser(res http.ResponseWriter, req *http.Request) {
	idUser, err := getID(req)
	if err != nil {
		h.log.Errorf("Get id user: %v", err)
		errMsg(res, h.templates)
		return
	}
	reqUser := &calendarapi.GetEventsForUserRequest{User: &calendarapi.UserRequest{Id: int64(idUser)},
		DaysAhead: &calendarapi.DaysAheadRequest{Days: 7, Date: timestamppb.New(time.Now())}}

	events, err := h.grpcclient.GetEventsForUser(req.Context(), reqUser)
	if err != nil {
		h.log.Errorf("displayListEventsForUser: GetEventsForUser: %v", err)
		errMsg(res, h.templates)
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
		h.log.Errorf("displayListEventsForUser: %v", err)
		errMsg(res, h.templates)
		return
	}
}

func (h *Handler) DeleteEvent(res http.ResponseWriter, req *http.Request) {
	idEvent, err := getID(req)
	if err != nil {
		h.log.Errorf("Get id event: %v", err)
		errMsg(res, h.templates)
		return
	}
	reqEvent := &calendarapi.EventRequest{Id: int64(idEvent)}
	if _, err := h.grpcclient.DeleteEvent(req.Context(), reqEvent); err != nil {
		h.log.Errorf("Delete event: %v", err)
		errMsg(res, h.templates)
		return
	}

	text := resp{Title: "Delete event", Body: struct{ Text string }{"Event deleted successfully"}}
	if err := h.templates.Execute("text.html", res, text); err != nil {
		h.log.Errorf("Delete event: %v", err)
		errMsg(res, h.templates)
		return
	}
}
