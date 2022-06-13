package internalhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/cmd/calendar/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	memorystorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/memory"
	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
)

const urlUsers = "http://localhost:8080/users"
const urlEvents = "http://localhost:8080/events"

var confDBTest = config.DBConf{"postgres", 5432, "event-db", "event-db", "event-db"}

func initServerHTTP(st storage.Storage) *Server {
	conf := config.HTTPConf{"localhost", 8080}
	server := New(&conf, st)

	return server
}

func TestHandleUserMemory(t *testing.T) {
	st := memorystorage.New()
	server := initServerHTTP(st)
	testHandleUser(t, server)
}

func TestHandleEventMemory(t *testing.T) {
	st := memorystorage.New()
	server := initServerHTTP(st)
	testHandleEvent(t, server)
}

func TestHandleUserSQL(t *testing.T) {
	st, err := sqlstorage.New(ctx, &confDBTest)
	if err != nil {
		t.Fatalf("db: %v", err)
	}
	defer st.Close()
	defer st.DropSchemaDB(ctx)

	server := initServerHTTP(st)
	testHandleUser(t, server)
}

func TestHandleEventSQL(t *testing.T) {
	st, err := sqlstorage.New(ctx, &confDBTest)
	if err != nil {
		t.Fatalf("db: %v", err)
	}
	defer st.Close()
	defer st.DropSchemaDB(ctx)

	server := initServerHTTP(st)
	testHandleEvent(t, server)
}

func testHandleUser(t *testing.T, server *Server) {

	users := []storage.User{
		{Name: "Howard Mendoza", Email: "Howard.Mendoza@mail.com"},
		{Name: "Brian Olson", Email: "B.Olson@gmail.com"},
		{Name: "Clarence Olson", Email: "Clarence.Olson@yandex.ru"},
	}

	// Create users
	for i := range users {
		testHandleCreateUser(t, server, &users[i], http.StatusCreated)
		users[i].ID = i + 1
	}

	// Get users
	for i := range users {
		user := testHandleGetUser(t, server, users[i].ID, http.StatusOK)
		if user != nil {
			if user.Name != users[i].Name {
				t.Errorf("mismatch name user: %s, %s", user.Name, users[i].Name)
			}
			if user.Email != users[i].Email {
				t.Errorf("mismatch email user: %s, %s", user.Email, users[i].Email)
			}
		}
	}

	// Update user email and get user
	for i := range users {
		users[i].Email = strings.ToLower(users[i].Email)
		testHandleUpdateUser(t, server, &users[i], http.StatusOK)

		user := testHandleGetUser(t, server, users[i].ID, http.StatusOK)
		if user != nil {
			if user.Email != users[i].Email {
				t.Errorf("mismatch email user: %s, %s", user.Email, users[i].Email)
			}
		}
	}

	// Delete all users and trying get, update, delete user that doesn't exist
	for i := range users {
		testHandleDeleteUser(t, server, users[i].ID, http.StatusOK)

		testHandleGetUser(t, server, users[i].ID, http.StatusBadRequest)
		testHandleUpdateUser(t, server, &users[i], http.StatusBadRequest)
		testHandleDeleteUser(t, server, users[i].ID, http.StatusBadRequest)
	}
}

func testHandleEvent(t *testing.T, server *Server) {

	user := storage.User{Name: "Bob", Email: "bobi@mail.com", Events: make([]storage.Event, 0)}
	user.ID = 1

	testHandleCreateUser(t, server, &user, http.StatusCreated)

	events := []storage.Event{
		{Title: "Bob's Birthday", Description: "Birthday February 24, 1993. Party in nature.",
			StartTime: time.Date(2022, time.February, 27, 10, 0, 0, 0, time.UTC),
			StopTime:  time.Date(2022, time.February, 27, 23, 0, 0, 0, time.UTC), UserID: user.ID},
		{Title: "Alis's Birthday", Description: "Birthday April 12, 1996. House party",
			StartTime: time.Date(2022, time.April, 13, 19, 0, 0, 0, time.UTC),
			StopTime:  time.Date(2022, time.April, 13, 21, 0, 0, 0, time.UTC), UserID: user.ID},
		{Title: "Jim's Birthday", Description: "Birthday August 15, 1994. Party at the restaurant",
			StartTime: time.Date(2022, time.August, 17, 16, 0, 0, 0, time.UTC),
			StopTime:  time.Date(2022, time.August, 17, 19, 0, 0, 0, time.UTC), UserID: user.ID},
		{Title: "Bill's Birthday", Description: "Birthday November 6, 1990. Party at the club",
			StartTime: time.Date(2022, time.November, 7, 11, 0, 0, 0, time.UTC),
			StopTime:  time.Date(2022, time.November, 7, 12, 0, 0, 0, time.UTC), UserID: user.ID},
	}

	// Create events
	for i := range events {
		testHandleCreateEvent(t, server, &events[i], http.StatusCreated)
		events[i].ID = i + 1
	}

	// Get events
	for i := range events {
		event := testHandleGetEvent(t, server, events[i].ID, http.StatusOK)
		if event != nil {
			if event.Title != events[i].Title {
				t.Errorf("mismatch title event: %s, %s", event.Title, events[i].Title)
			}
			if event.Description != events[i].Description {
				t.Errorf("mismatch description event: %s, %s", event.Description, events[i].Description)
			}
		}
	}

	// Update event title and get event
	for i := range events {
		events[i].Title = strings.ToUpper(events[i].Title)
		testHandleUpdateEvent(t, server, &events[i], http.StatusOK)

		event := testHandleGetEvent(t, server, events[i].ID, http.StatusOK)
		if event != nil {
			if event.Title != events[i].Title {
				t.Errorf("mismatch title event: %s, %s", event.Title, events[i].Title)
			}
		}
	}

	// Delete all event without last and trying get, update, delete event that doesn't exist
	for i := 0; i < len(events)-1; i++ {
		testHandleDeleteEvent(t, server, events[i].ID, http.StatusOK)

		testHandleGetEvent(t, server, events[i].ID, http.StatusBadRequest)
		testHandleUpdateEvent(t, server, &events[i], http.StatusBadRequest)
		testHandleDeleteEvent(t, server, events[i].ID, http.StatusBadRequest)
	}

	testHandleDeleteEvent(t, server, len(events), http.StatusOK)
	testHandleDeleteUser(t, server, user.ID, http.StatusOK)
	testHandleDeleteEvent(t, server, len(events), http.StatusBadRequest)
}

func testHandleCreateUser(t *testing.T, server *Server, user *storage.User, status int) {
	res := httptest.NewRecorder()

	dataJsonUser, err := json.Marshal(*user)
	if err != nil {
		t.Fatalf("HandleCreateUser: cant marshal JSON: %v", err)
	}
	req, err := http.NewRequest("POST", urlUsers, bytes.NewReader(dataJsonUser))
	if err != nil {
		t.Fatalf("HandleCreateUser: create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	server.ServeHTTP(res, req)

	if res.Code != status {
		t.Errorf("HandleCreateUser: response code is %d; want: %d", res.Code, status)
	}
}

func testHandleGetUser(t *testing.T, server *Server, id int, status int) *storage.User {
	res := httptest.NewRecorder()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?id=%d", urlUsers, id), nil)
	if err != nil {
		t.Fatalf("HandleGetUser: create request: %v", err)
	}

	server.ServeHTTP(res, req)

	if res.Code != status {
		t.Errorf("HandleGetUser: response code is %d; want: %d", res.Code, status)
	}
	if res.Code != http.StatusOK {
		return nil
	}

	var user storage.User
	if err := json.Unmarshal(res.Body.Bytes(), &user); err != nil {
		t.Fatalf("HandleGetUser: cant unmarshal JSON: %v", err)
	}

	return &user
}

func testHandleUpdateUser(t *testing.T, server *Server, user *storage.User, status int) {
	res := httptest.NewRecorder()

	dataJsonUser, err := json.Marshal(*user)
	if err != nil {
		t.Fatalf("HandleUpdateUser: cant marshal JSON: %v", err)
	}
	req, err := http.NewRequest("PUT", urlUsers, bytes.NewReader(dataJsonUser))
	if err != nil {
		t.Fatalf("HandleUpdateUser: create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	server.ServeHTTP(res, req)

	if res.Code != status {
		t.Errorf("HandleUpdateUser: response code is %d; want: %d", res.Code, status)
	}
}

func testHandleDeleteUser(t *testing.T, server *Server, id int, status int) {
	res := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s?id=%d", urlUsers, id), nil)
	if err != nil {
		t.Fatalf("HandleDeleteUser: create request %v", err)
	}

	server.ServeHTTP(res, req)

	if res.Code != status {
		t.Errorf("HandleDeleteUser: response code is %d; want: %d", res.Code, status)
	}
}

func testHandleCreateEvent(t *testing.T, server *Server, event *storage.Event, status int) {
	res := httptest.NewRecorder()

	dataJsonEvent, err := json.Marshal(*event)
	if err != nil {
		t.Fatalf("HandleCreateEvent: cant marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", urlEvents, bytes.NewReader(dataJsonEvent))
	if err != nil {
		t.Fatalf("HandleCreateEvent: create request %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	server.ServeHTTP(res, req)

	if res.Code != status {
		t.Errorf("HandleCreateEvent: response code is %d; want: %d", res.Code, status)
	}
}

func testHandleGetEvent(t *testing.T, server *Server, id int, status int) *storage.Event {
	res := httptest.NewRecorder()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?id=%d", urlEvents, id), nil)
	if err != nil {
		t.Fatalf("HandleGetEvent: create request: %v", err)
	}

	server.ServeHTTP(res, req)

	if res.Code != status {
		t.Errorf("HandleGetEvent: response code is %d; want: %d", res.Code, status)
	}
	if res.Code != http.StatusOK {
		return nil
	}

	var event storage.Event
	if err := json.Unmarshal(res.Body.Bytes(), &event); err != nil {
		t.Errorf("HandleGetEvent: cant unmarshal JSON: %v", err)
	}

	return &event
}

func testHandleUpdateEvent(t *testing.T, server *Server, event *storage.Event, status int) {
	res := httptest.NewRecorder()

	dataJsonEvent, err := json.Marshal(*event)
	if err != nil {
		t.Fatalf("HandleUpdateEvent: cant marshal JSON: %v", err)
	}
	req, err := http.NewRequest("PUT", urlEvents, bytes.NewReader(dataJsonEvent))
	if err != nil {
		t.Fatalf("HandleUpdateEvent: create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	server.ServeHTTP(res, req)

	if res.Code != status {
		t.Errorf("HandleUpdateEvent: response code is %d; want: %d", res.Code, status)
	}
}

func testHandleDeleteEvent(t *testing.T, server *Server, id int, status int) {
	res := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s?id=%d", urlEvents, id), nil)
	if err != nil {
		t.Fatalf("HandleDeleteEvent: create request: %v", err)
	}

	server.ServeHTTP(res, req)

	if res.Code != status {
		t.Errorf("HandleDeleteEvent: response code is %d; want: %d", res.Code, status)
	}
}
