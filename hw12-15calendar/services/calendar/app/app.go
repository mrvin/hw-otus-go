package app

import (
	"context"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

type App struct {
	storage storage.Storage
}

func New(storage storage.Storage) *App {
	return &App{storage}
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) GetEvent(ctx context.Context, id int) (*storage.Event, error) {
	return a.storage.GetEvent(ctx, id)
}

/*
func (a *App) GetAllEvents(ctx context.Context) ([]storage.Event, error) {
	return a.storage.GetAllEvents(ctx)
}
*/

func (a *App) GetEventsForUser(ctx context.Context, id int) ([]storage.Event, error) {
	return a.storage.GetEventsForUser(ctx, id)
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) CreateUser(ctx context.Context, user *storage.User) error {
	return a.storage.CreateUser(ctx, user)
}

func (a *App) GetUser(ctx context.Context, id int) (*storage.User, error) {
	return a.storage.GetUser(ctx, id)
}

func (a *App) GetAllUsers(ctx context.Context) ([]storage.User, error) {
	return a.storage.GetAllUsers(ctx)
}

func (a *App) UpdateUser(ctx context.Context, user *storage.User) error {
	return a.storage.UpdateUser(ctx, user)
}

func (a *App) DeleteUser(ctx context.Context, id int) error {
	return a.storage.DeleteUser(ctx, id)
}
