package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/event"
	"github.com/jasonsoft/starter/pkg/event/proto"
)

type EventServer struct {
	config       config.Configuration
	eventService event.EventServicer
}

func NewEventServer(cfg config.Configuration, eventService event.EventServicer) proto.EventServiceServer {
	return &EventServer{
		config:       cfg,
		eventService: eventService,
	}
}

// GetEvents returns all events
func (s *EventServer) GetEvents(ctx context.Context, request *empty.Empty) (*proto.GetEventsResponse, error) {
	events, err := s.eventService.Events(ctx)
	if err != nil {
		return nil, err
	}

	result := proto.GetEventsResponse{
		Data: convertToGRPCEvents(events),
	}

	return &result, nil
}

// UpdatePublishStatus update event's publishstatus
func (s *EventServer) UpdatePublishStatus(ctx context.Context, request *proto.UpdatePublishStatusRequest) (*empty.Empty, error) {

	req := event.UpdateEventStatusRequest{
		EventID:         request.EventId,
		TransID:         request.TransId,
		PublishedStatus: event.PublishedStatus(request.PublishedStatus),
	}

	err := s.eventService.UpdatePublishStatus(ctx, req)
	return &empty.Empty{}, err
}

func convertToGRPCEvents(events []*event.Event) []*proto.Event {
	result := []*proto.Event{}

	for _, evt := range events {
		createdAt, _ := ptypes.TimestampProto(evt.CreatedAt)

		target := proto.Event{
			Id:              evt.ID,
			Title:           evt.Title,
			Description:     evt.Description,
			PublishedStatus: proto.PublishedStatus(evt.PublishedStatus),
			CreatedAt:       createdAt,
		}

		result = append(result, &target)
	}

	return result
}
