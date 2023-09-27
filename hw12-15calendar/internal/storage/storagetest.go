package storage

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"
	"time"
)

var users = []User{
	{Name: "Bob", Email: "bobi@mail.com", Events: make([]Event, 0)},
	{Name: "Alis", Email: "alisia.jones@gmail.com", Events: make([]Event, 0)},
	{Name: "Jim", Email: "jimihendrix@yandex.ru", Events: make([]Event, 0)},
}

var events = []Event{
	{Title: "Bob's Birthday", Description: "Birthday February 24, 1993. Party in nature.",
		StartTime: time.Date(2022, time.February, 27, 10, 0, 0, 0, time.UTC),
		StopTime:  time.Date(2022, time.February, 27, 23, 0, 0, 0, time.UTC)},
	{Title: "Alis's Birthday", Description: "Birthday April 12, 1996. House party",
		StartTime: time.Date(2022, time.April, 13, 19, 0, 0, 0, time.UTC),
		StopTime:  time.Date(2022, time.April, 13, 21, 0, 0, 0, time.UTC)},
	{Title: "Jim's Birthday", Description: "Birthday August 15, 1994. Party at the restaurant",
		StartTime: time.Date(2022, time.August, 17, 16, 0, 0, 0, time.UTC),
		StopTime:  time.Date(2022, time.August, 17, 19, 0, 0, 0, time.UTC)},
	{Title: "Bill's Birthday", Description: "Birthday November 6, 1990. Party at the club",
		StartTime: time.Date(2022, time.November, 7, 11, 0, 0, 0, time.UTC),
		StopTime:  time.Date(2022, time.November, 7, 12, 0, 0, 0, time.UTC)},
}

type byIDEvent []Event

func (x byIDEvent) Len() int           { return len(x) }
func (x byIDEvent) Less(i, j int) bool { return x[i].ID < x[j].ID }
func (x byIDEvent) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func TestUserCRUD(ctx context.Context, t *testing.T, st Storage) {
	// Create users
	for i := range users {
		if err := st.CreateUser(ctx, &users[i]); err != nil {
			t.Errorf("CreateUser: %v", err)
		}
		if users[i].ID == 0 {
			t.Errorf("CreateUser: can't get ID")
		}
	}

	// Get users without events
	for i := range users {
		user, err := st.GetUser(ctx, users[i].ID)
		if err != nil {
			t.Errorf("GetUser(id = %d): %v", users[i].ID, err)
		}

		if !reflect.DeepEqual(*user, users[i]) {
			t.Errorf("GetUser(id = %d):\n\thave: %v\n\twant: %v", users[i].ID, *user, users[i])
		}
	}

	_, err := st.ListUsers(ctx)
	if err != nil {
		t.Errorf("GetAllUsers: %v", err)
	}

	// Update user name
	users[0].Name = "Bill"
	if err := st.UpdateUser(ctx, &users[0]); err != nil {
		t.Errorf("UpdateUser(id = %d): %v", users[0].ID, err)
	}
	user, err := st.GetUser(ctx, users[0].ID)
	if err != nil {
		t.Errorf("UpdateUser: get user with id = %d: %v", users[0].ID, err)
	}
	if user.Name != users[0].Name {
		t.Errorf("UpdateUser:(id = %d):\n\thave: %v\n\twant: %v", users[0].ID, *user, users[0])
	}

	// Delete all users
	for _, user := range users {
		if err := st.DeleteUser(ctx, user.ID); err != nil {
			t.Errorf("DeleteUser: %v", err)
		}
	}

	// Trying get, update, delete user that doesn't exist
	_, err = st.GetUser(ctx, users[0].ID)
	if !errors.Is(err, ErrNoUser) {
		t.Errorf("GetUser(id = %d): %v", users[0].ID, err)
	}
	err = st.UpdateUser(ctx, &users[1])
	if !errors.Is(err, ErrNoUser) {
		t.Errorf("UpdateUser(id = %d): %v", users[1].ID, err)
	}
	err = st.DeleteUser(ctx, users[2].ID)
	if !errors.Is(err, ErrNoUser) {
		t.Errorf("DeleteUser(id = %d): %v", users[2].ID, err)
	}
}

func TestEventCRUD(ctx context.Context, t *testing.T, st Storage) { //nolint:funlen,gocognit,cyclop
	// Create users
	for i := range users {
		if err := st.CreateUser(ctx, &users[i]); err != nil {
			t.Errorf("CreateUser: %v", err)
		}
		if users[i].ID == 0 {
			t.Errorf("CreateUser: can't get ID")
		}
	}

	// Create events for users
	for i := range users {
		for j := range events {
			if i != j {
				events[j].UserID = users[i].ID

				if err := st.CreateEvent(ctx, &events[j]); err != nil {
					t.Errorf("CreateEvent: %v", err)
				}
				if events[j].ID == 0 {
					t.Errorf("CreateEvent: can't get ID")
				}

				users[i].Events = append(users[i].Events, events[j])
			}
		}
	}

	// Get users with events
	for i := range users {
		user, err := st.GetUser(ctx, users[i].ID)
		if err != nil {
			t.Errorf("GetUser(id = %d): %v", users[i].ID, err)
		}

		sort.Sort(byIDEvent(user.Events))
		cmpUsers(t, user, &users[i])
	}

	// Update event
	events[0].Description = "Cool Birthday"
	if err := st.UpdateEvent(ctx, &events[0]); err != nil {
		t.Errorf("UpdateEvent(id = %d): %v", events[0].ID, err)
	}

	event, err := st.GetEvent(ctx, events[0].ID)
	if err != nil {
		t.Errorf("GetEvent(id = %d): %v", events[0].ID, err)
	}
	cmpEvent(t, event, &events[0])

	_, err = st.ListEvents(ctx)
	if err != nil {
		t.Errorf("GetAllEvents: %v", err)
	}

	// Delete all events
	for _, user := range users {
		for _, event := range user.Events {
			if err := st.DeleteEvent(ctx, event.ID); err != nil {
				t.Errorf("DeleteEvent: %v", err)
			}
		}
	}

	// Delete all users
	for _, user := range users {
		if err := st.DeleteUser(ctx, user.ID); err != nil {
			t.Errorf("DeleteUser: %v", err)
		}
	}

	// Trying get, update, delete event that doesn't exist
	_, err = st.GetEvent(ctx, events[0].ID)
	if !errors.Is(err, ErrNoEvent) {
		t.Errorf("GetEvent(id = %d): %v", events[0].ID, err)
	}
	err = st.UpdateEvent(ctx, &events[1])
	if !errors.Is(err, ErrNoEvent) {
		t.Errorf("UpdateEvent(id = %d): %v", events[1].ID, err)
	}
	err = st.DeleteEvent(ctx, events[2].ID)
	if !errors.Is(err, ErrNoEvent) {
		t.Errorf("DeleteEvent(id = %d): %v", events[2].ID, err)
	}
}

func cmpUsers(t *testing.T, u1, u2 *User) {
	if u1.ID != u2.ID {
		t.Errorf("mismatch id user: %d, %d", u1.ID, u2.ID)
	}
	if u1.Name != u2.Name {
		t.Errorf("mismatch name user: %s, %s", u1.Name, u2.Name)
	}
	if u1.Email != u2.Email {
		t.Errorf("mismatch email user: %s, %s", u1.Email, u2.Email)
	}

	for i := range u1.Events {
		cmpEvent(t, &u1.Events[i], &u2.Events[i])
	}
}

func cmpEvent(t *testing.T, e1, e2 *Event) {
	if e1.ID != e2.ID {
		t.Errorf("mismatch id event: %d, %d", e1.ID, e2.ID)
	}
	if e1.Title != e2.Title {
		t.Errorf("mismatch title event: %s, %s", e1.Title, e2.Title)
	}
	if e1.Description != e2.Description {
		t.Errorf("mismatch description event: %s, %s", e1.Description, e2.Description)
	}
	if !e1.StartTime.Equal(e2.StartTime) {
		t.Errorf("mismatch start time event: %v, %v", e1.StartTime, e2.StartTime)
	}
	if !e1.StopTime.Equal(e2.StopTime) {
		t.Errorf("mismatch stop time event: %v, %v", e1.StopTime, e2.StopTime)
	}
	if e1.UserID != e2.UserID {
		t.Errorf("mismatch user id for event: %d, %d", e1.UserID, e2.UserID)
	}
}
