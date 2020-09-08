package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/domain"
	"github.com/jasonsoft/starter/pkg/event/proto"
)

// EventServer handles all event business logic
type EventServer struct {
	config       config.Configuration
	eventService domain.EventServicer
}

// NewEventServer create an instance of EventServer
func NewEventServer(cfg config.Configuration, eventService domain.EventServicer) proto.EventServiceServer {
	return &EventServer{
		config:       cfg,
		eventService: eventService,
	}
}

// GetEvents returns all events
func (s *EventServer) GetEvents(ctx context.Context, request *proto.GetEventsRequest) (*proto.GetEventsResponse, error) {
	logger := log.FromContext(ctx)
	logger.Debug("grpc: begin GetEvent fn")

	events, err := s.eventService.Events(ctx, domain.FindEventOptions{
		ID:    request.Id,
		Title: request.Title,
	})
	if err != nil {
		return nil, err
	}

	data, err := eventsToGRPC(events)
	if err != nil {
		return nil, nil
	}
	result := proto.GetEventsResponse{
		Events: data,
	}

	return &result, nil
}

// UpdatePublishStatus update event's publishstatus
func (s *EventServer) UpdatePublishStatus(ctx context.Context, request *proto.UpdatePublishStatusRequest) (*empty.Empty, error) {

	req := domain.UpdateEventStatusRequest{
		EventID:         request.EventId,
		TransID:         request.TransId,
		PublishedStatus: domain.PublishedStatus(request.PublishedStatus),
	}

	err := s.eventService.UpdatePublishStatus(ctx, req)
	return &empty.Empty{}, err
}
