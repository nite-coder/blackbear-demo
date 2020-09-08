package grpc

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/jasonsoft/starter/pkg/domain"
	"github.com/jasonsoft/starter/pkg/event/proto"
)

func eventToGRPC(event *domain.Event) (*proto.Event, error) {
	if event == nil {
		return nil, nil
	}

	createdAt, err := ptypes.TimestampProto(event.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := ptypes.TimestampProto(event.UpdatedAt)
	if err != nil {
		return nil, err
	}

	result := proto.Event{
		Id:              event.ID,
		Title:           event.Title,
		Description:     event.Description,
		PublishedStatus: proto.PublishedStatus(event.PublishedStatus),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
	return &result, nil
}

func eventsToGRPC(events []domain.Event) ([]*proto.Event, error) {
	result := []*proto.Event{}

	for _, evt := range events {
		event, err := eventToGRPC(&evt)
		if err != nil {
			return nil, err
		}
		result = append(result, event)
	}

	return result, nil
}
