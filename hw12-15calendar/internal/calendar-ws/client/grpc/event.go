package grpcclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-api"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) CreateEvent(
	ctx context.Context, token, title, description string,
	startTime, stopTime time.Time,
) (int64, error) {

	event := &calendarapi.CreateEventRequest{
		AccessToken: token,
		Title:       title,
		Description: description,
		StartTime:   timestamppb.New(startTime),
		StopTime:    timestamppb.New(stopTime),
	}
	response, err := c.eventService.CreateEvent(ctx, event)
	if err != nil {
		return 0, fmt.Errorf("gRPC: %w", err)
	}
	slog.Debug("Added event",
		slog.Int64("id", response.Id),
		slog.String("title", event.Title),
		slog.String("Description", event.Description),
	)

	return response.Id, nil
}

func (c *Client) GetEvent(ctx context.Context, token string, id int64) (*storage.Event, error) {
	reqEvent := &calendarapi.GetEventRequest{
		AccessToken: token,
		Id:          id,
	}
	event, err := c.eventService.GetEvent(ctx, reqEvent)
	if err != nil {
		return nil, fmt.Errorf("gRPC: %w", err)
	}

	return &storage.Event{
		ID:          event.Id,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime.AsTime(),
		StopTime:    event.StopTime.AsTime(),
	}, nil
}

func (c *Client) UpdateEvent(
	ctx context.Context, token, title, description string,
	startTime, stopTime time.Time,
	userID uuid.UUID) error {

	event := &calendarapi.UpdateEventRequest{
		AccessToken: token,
		Title:       title,
		Description: description,
		StartTime:   timestamppb.New(startTime),
		StopTime:    timestamppb.New(stopTime),
		UserID:      &calendarapi.UUID{Value: userID.String()},
	}
	if _, err := c.eventService.UpdateEvent(ctx, event); err != nil {
		return fmt.Errorf("gRPC: %w", err)
	}

	return nil
}

func (c *Client) DeleteEvent(ctx context.Context, token string, id int64) error {
	reqEvent := &calendarapi.DeleteEventRequest{
		AccessToken: token,
		Id:          id,
	}
	if _, err := c.eventService.DeleteEvent(ctx, reqEvent); err != nil {
		return fmt.Errorf("gRPC: %w", err)

	}

	return nil
}

func (c *Client) ListEventsForUser(ctx context.Context, token string, days int) ([]storage.Event, error) {
	reqUser := &calendarapi.ListEventsForUserRequest{
		AccessToken: token,
		Days:        int32(days),
		Date:        timestamppb.New(time.Now()),
	}

	events, err := c.eventService.ListEventsForUser(ctx, reqUser)
	if err != nil {
		return nil, fmt.Errorf("gRPC: %w", err)
	}
	fmt.Printf("Events: %v\n", events)
	listEvents := make([]storage.Event, 0, len(events.Events))
	for _, event := range events.Events {
		eventUserIDUUID, err := uuid.Parse(event.GetUserID().GetValue())
		if err != nil {
			return nil, fmt.Errorf("gRPC: %w", err)
		}
		listEvents = append(listEvents, storage.Event{
			ID:          event.Id,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime.AsTime(),
			StopTime:    event.StopTime.AsTime(),
			UserID:      eventUserIDUUID,
		})
	}

	return listEvents, nil
}
