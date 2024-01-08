package grpcserver

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/calendar-api"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateEvent(ctx context.Context, pbEvent *calendarapi.CreateEventRequest) (*calendarapi.CreateEventResponse, error) {
	if err := pbEvent.StartTime.CheckValid(); err != nil {
		err = fmt.Errorf("incorrect value StartTime: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	if err := pbEvent.StopTime.CheckValid(); err != nil {
		err = fmt.Errorf("incorrect value StopTime: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	userUUID, err := uuid.Parse(pbEvent.GetUserID().GetValue())
	if err != nil {
		err = fmt.Errorf("parse uuid: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	event := storage.Event{
		ID:          0,
		Title:       pbEvent.GetTitle(),
		Description: pbEvent.GetDescription(),
		StartTime:   pbEvent.StartTime.AsTime(),
		StopTime:    pbEvent.StopTime.AsTime(),
		UserID:      userUUID,
	}

	id, err := s.eventService.CreateEvent(ctx, &event)
	if err != nil {
		err = fmt.Errorf("create event: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &calendarapi.CreateEventResponse{Id: id}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *calendarapi.GetEventRequest) (*calendarapi.EventResponse, error) {
	event, err := s.eventService.GetEvent(ctx, req.GetId())
	if err != nil {
		err := fmt.Errorf("get event: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	userUUID := calendarapi.UUID{Value: event.UserID.String()}
	return &calendarapi.EventResponse{
		Id:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   timestamppb.New(event.StartTime),
		StopTime:    timestamppb.New(event.StopTime),
		UserID:      &userUUID,
	}, nil
}

func (s *Server) ListEventsForUser(ctx context.Context, req *calendarapi.ListEventsForUserRequest) (*calendarapi.ListEventsResponse, error) {
	if err := req.Date.CheckValid(); err != nil {
		return nil, fmt.Errorf("incorrect value date: %w", err)
	}
	date := req.Date.AsTime()

	userName := GetUserName(ctx)
	if userName == "" {
		err := fmt.Errorf("user name is empty")
		slog.Error(err.Error())
		return nil, err
	}
	events, err := s.eventService.ListEventsForUser(ctx, userName, date, int(req.Days))
	if err != nil {
		err := fmt.Errorf("get events for user: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	pbEvents := make([]*calendarapi.EventResponse, len(events), len(events))
	for i, event := range events {
		userUUID := calendarapi.UUID{Value: event.UserID.String()}
		pbEvents[i] = &calendarapi.EventResponse{
			Id:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   timestamppb.New(event.StartTime),
			StopTime:    timestamppb.New(event.StopTime),
			UserID:      &userUUID,
		}
	}

	return &calendarapi.ListEventsResponse{Events: pbEvents}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, pbEvent *calendarapi.UpdateEventRequest) (*emptypb.Empty, error) {
	if err := pbEvent.StartTime.CheckValid(); err != nil {
		err = fmt.Errorf("incorrect value StartTime: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	if err := pbEvent.StopTime.CheckValid(); err != nil {
		err = fmt.Errorf("incorrect value StopTime: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	userUUID, err := uuid.Parse(pbEvent.GetUserID().GetValue())
	if err != nil {
		err = fmt.Errorf("parse uuid: %w", err)
		slog.Error(err.Error())
		return nil, err
	}
	event := storage.Event{
		ID:          0,
		Title:       pbEvent.GetTitle(),
		Description: pbEvent.GetDescription(),
		StartTime:   pbEvent.StartTime.AsTime(),
		StopTime:    pbEvent.StopTime.AsTime(),
		UserID:      userUUID,
	}

	if err := s.eventService.UpdateEvent(ctx, &event); err != nil {
		err := fmt.Errorf("update event: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *calendarapi.DeleteEventRequest) (*emptypb.Empty, error) {
	if err := s.eventService.DeleteEvent(ctx, req.GetId()); err != nil {
		err := fmt.Errorf("delete event: %w", err)
		slog.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
