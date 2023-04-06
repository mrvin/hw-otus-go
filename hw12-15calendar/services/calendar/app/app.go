package app

import (
	"context"
	"errors"
	"time"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var ErrStopTimeBeforeStartTime = errors.New("event ends before starts")

type App struct {
	storage storage.Storage
	tr      trace.Tracer // Think about it.
}

func New(storage storage.Storage) *App {
	return &App{storage, otel.GetTracerProvider().Tracer("Storage")}
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	if event.StopTime.Before(event.StartTime) {
		return ErrStopTimeBeforeStartTime
	}

	cctx, sp := a.tr.Start(ctx, "CreateEvent")
	defer sp.End()

	return a.storage.CreateEvent(cctx, event)
}

func (a *App) GetEvent(ctx context.Context, id int) (*storage.Event, error) {
	cctx, sp := a.tr.Start(ctx, "GetEvent")
	defer sp.End()

	return a.storage.GetEvent(cctx, id)
}

/*
func (a *App) GetAllEvents(ctx context.Context) ([]storage.Event, error) {
	return a.storage.GetAllEvents(ctx)
}
*/

// TODO: implement at the database level.
func (a *App) GetEventsForUser(ctx context.Context, id int, startPeriod time.Time, days int) ([]storage.Event, error) {
	cctx, sp := a.tr.Start(ctx, "GetEventsForUser")
	defer sp.End()

	events, err := a.storage.GetEventsForUser(cctx, id)
	if err != nil {
		return nil, err
	}
	if days == 0 {
		return events, nil
	}
	stopPeriod := startPeriod.AddDate(0, 0, days)
	var eventsFromInterval []storage.Event
	for _, event := range events {
		if event.StartTime.After(startPeriod) && event.StartTime.Before(stopPeriod) ||
			event.StopTime.After(startPeriod) && event.StopTime.Before(stopPeriod) {
			eventsFromInterval = append(eventsFromInterval, event)
		}
	}

	return eventsFromInterval, nil
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) error {
	cctx, sp := a.tr.Start(ctx, "UpdateEvent")
	defer sp.End()

	return a.storage.UpdateEvent(cctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	cctx, sp := a.tr.Start(ctx, "DeleteEvent")
	defer sp.End()

	return a.storage.DeleteEvent(cctx, id)
}

func (a *App) CreateUser(ctx context.Context, user *storage.User) error {
	cctx, sp := a.tr.Start(ctx, "CreateUser")
	defer sp.End()

	return a.storage.CreateUser(cctx, user)
}

func (a *App) GetUser(ctx context.Context, id int) (*storage.User, error) {
	cctx, sp := a.tr.Start(ctx, "GetUser")
	defer sp.End()

	return a.storage.GetUser(cctx, id)
}

func (a *App) GetAllUsers(ctx context.Context) ([]storage.User, error) {
	cctx, sp := a.tr.Start(ctx, "GetAllUsers")
	defer sp.End()

	return a.storage.GetAllUsers(cctx)
}

func (a *App) UpdateUser(ctx context.Context, user *storage.User) error {
	cctx, sp := a.tr.Start(ctx, "UpdateUser")
	defer sp.End()

	return a.storage.UpdateUser(cctx, user)
}

func (a *App) DeleteUser(ctx context.Context, id int) error {
	cctx, sp := a.tr.Start(ctx, "DeleteUser")
	defer sp.End()

	return a.storage.DeleteUser(cctx, id)
}
