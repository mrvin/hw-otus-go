package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendarapi"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func displayFormEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	idUser, err := getID(req)
	if err != nil {
		log.Printf("Get id user: %v", err)
		errMsg(res, server.templates)
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
	if err := executeTemp("form-event.html", res, data, server.templates); err != nil {
		log.Printf("displayFormEvent: %v", err)
		errMsg(res, server.templates)
		return
	}
}

func createEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	if err := req.ParseForm(); err != nil {
		log.Printf("createEvent: %v", err)
		errMsg(res, server.templates)
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
		fmt.Printf("get location: %v", err)
		errMsg(res, server.templates)
		return
	}
	idUser, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("convert id: %v", err)
		errMsg(res, server.templates)
		return
	}
	timeNow := time.Now().In(loc)
	tZone := timeNow.Format("Z07:00")
	startTimeGO, err := time.ParseInLocation(time.RFC3339, startTime+":00"+tZone, loc)
	if err != nil {
		log.Printf("сreateEvent: parse startTime: %v", err)
		errMsg(res, server.templates)
		return
	}

	stopTimeGO, err := time.ParseInLocation(time.RFC3339, stopTime+":00"+tZone, loc)
	if err != nil {
		log.Printf("сreateEvent: parse stopTime: %v", err)
		errMsg(res, server.templates)
		return
	}

	starTimePB := timestamppb.New(startTimeGO)
	stopTimePB := timestamppb.New(stopTimeGO)
	event := &calendarapi.Event{Title: title, Description: description, StartTime: starTimePB, StopTime: stopTimePB, UserID: int64(idUser)}
	_, err = server.grpcclient.CreateEvent(req.Context(), event)
	if err != nil {
		log.Printf("сreateEvent: %v", err)
		errMsg(res, server.templates)
		return
	}

	text := resp{Title: "Create event", Body: struct{ Text string }{"Event created successfully"}}
	if err := executeTemp("text.html", res, text, server.templates); err != nil {
		log.Printf("createEvent: %v", err)
		errMsg(res, server.templates)
		return
	}
}

func displayEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	idEvent, err := getID(req)
	if err != nil {
		log.Printf("Get id event: %v", err)
		errMsg(res, server.templates)
		return
	}
	reqEvent := &calendarapi.EventRequest{Id: int64(idEvent)}
	event, err := server.grpcclient.GetEvent(req.Context(), reqEvent)
	if err != nil {
		log.Printf("displayEvent: %v", err)
		errMsg(res, server.templates)
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

	if err := executeTemp("event.html", res, dataEvent, server.templates); err != nil {
		log.Printf("displayEvent: %v", err)
		errMsg(res, server.templates)
		return
	}
}

func displayListEventsForUser(res http.ResponseWriter, req *http.Request, server *Server) {
	idUser, err := getID(req)
	if err != nil {
		log.Printf("Get id user: %v", err)
		errMsg(res, server.templates)
		return
	}
	reqUser := &calendarapi.GetEventsForUserRequest{User: &calendarapi.UserRequest{Id: int64(idUser)},
		DaysAhead: &calendarapi.DaysAheadRequest{Days: 7, Date: timestamppb.New(time.Now())}}

	events, err := server.grpcclient.GetEventsForUser(req.Context(), reqUser)
	if err != nil {
		log.Printf("displayListEventsForUser: GetEventsForUser: %v", err)
		errMsg(res, server.templates)
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

	if err := executeTemp("list-events.html", res, dataEvent, server.templates); err != nil {
		log.Printf("displayListEventsForUser: %v", err)
		errMsg(res, server.templates)
		return
	}
}

func deleteEvent(res http.ResponseWriter, req *http.Request, server *Server) {
	idEvent, err := getID(req)
	if err != nil {
		log.Printf("Get id event: %v", err)
		errMsg(res, server.templates)
		return
	}
	reqEvent := &calendarapi.EventRequest{Id: int64(idEvent)}
	if _, err := server.grpcclient.DeleteEvent(req.Context(), reqEvent); err != nil {
		log.Printf("Delete event: %v", err)
		errMsg(res, server.templates)
		return
	}

	text := resp{Title: "Delete event", Body: struct{ Text string }{"Event deleted successfully"}}
	if err := executeTemp("text.html", res, text, server.templates); err != nil {
		log.Printf("Delete event: %v", err)
		errMsg(res, server.templates)
		return
	}
}
