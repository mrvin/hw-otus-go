package sqlstorage

import (
	"testing"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

func TestStorageSQL(t *testing.T) {
	var st Storage
	conf := config.DBConf{"localhost", 5432, "event-db", "event-db", "event-db"}

	if err := st.Connect(nil, &conf); err != nil {
		t.Errorf("Connect: %v", err)
	}
	users := []storage.User{
		{Name: "Bob"},
		{Name: "Alis"},
		{Name: "Jim"},
	}
	for i, user := range users {
		if err := st.CreateUser(nil, &user); err != nil {
			t.Errorf("CreateUser: %v", err)
		}
		if user.ID == 0 {
			t.Errorf("CreateUser: can't get ID")
		}
		users[i].ID = user.ID
	}

	user, err := st.GetUser(nil, users[0].ID)
	if err != nil {
		t.Errorf("GetUser(id = %d): %v", users[0].ID, err)
	}
	if user.Name != users[0].Name {
		t.Errorf("GetUser: user name mismatch: %v", err)
	}

	users[0].Name = "Bill"
	if err := st.UpdateUser(nil, &users[0]); err != nil {
		t.Errorf("UpdateUser: %v", err)
	}

	if err := st.DeleteUser(nil, users[1].ID); err != nil {
		t.Errorf("DeleteUser: %v", err)
	}

	events := []storage.Event{
		{Title: "Bob's Birthday", Description: "P", StartTime: time.Date(1993, time.February, 27, 10, 0, 0, 0, time.UTC),
			StopTime: time.Date(1993, time.February, 27, 23, 0, 0, 0, time.UTC),
			UserID:   users[0].ID},
		{Title: "Alis's Birthday", Description: "R", UserID: users[2].ID},
		{Title: "Jim's Birthday", Description: "C", UserID: users[0].ID},
		{Title: "Bill's Birthday", Description: "S", UserID: users[2].ID},
	}

	for i, event := range events {
		if err := st.CreateEvent(nil, &event); err != nil {
			t.Errorf("CreateEvent: %v", err)
		}
		if event.ID == 0 {
			t.Errorf("CreateEvent: can't get ID")
		}
		events[i].ID = event.ID
	}

	event, err := st.GetEvent(nil, events[0].ID)
	if err != nil {
		t.Errorf("GetEvent(id = %d): %v", events[0].ID, err)
	}
	if event.Title != events[0].Title || event.Description != events[0].Description {
		t.Errorf("GetEvent: event mismatch: %v", err)
	}

	events[0].Description = "Cool Birthday"
	if err := st.UpdateEvent(nil, &events[0]); err != nil {
		t.Errorf("UpdateEvent: %v", err)
	}

	if err := st.DeleteEvent(nil, events[1].ID); err != nil {
		t.Errorf("DeleteEvent: %v", err)
	}
	if err := st.DeleteEvent(nil, events[3].ID); err != nil {
		t.Errorf("DeleteEvent: %v", err)
	}

	if err := st.DeleteUser(nil, users[0].ID); err != nil {
		t.Errorf("DeleteUser: %v", err)
	}

	if err := st.DeleteUser(nil, users[2].ID); err != nil {
		t.Errorf("DeleteUser: %v", err)
	}

	st.Close(nil)
}
