package gql

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jasonsoft/log/v2"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *queryResolver) GetEvents(ctx context.Context) ([]*Event, error) {
	logger := log.FromContext(ctx)

	logger.Debug("GetEvents 1")

	resp, err := r.eventClient.GetEvents(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("eventClient call failed, name: %s: %w", "GetEvents", err)
	}

	logger.Debug("GetEvents 2")
	result := []*Event{}

	for _, d := range resp.Data {
		result = append(result, grpcEventToGQLEvent(d))
		logger.Debug("GetEvents 2.1")
	}

	logger.Debug("GetEvents 3")
	return result, nil
}

func (r *queryResolver) GetWallet(ctx context.Context) (*Wallet, error) {
	panic("not implemented")
}

func grpcEventToGQLEvent(source *eventProto.Event) *Event {

	createdAt, _ := ptypes.Timestamp(source.CreatedAt)
	result := Event{
		ID:          source.Id,
		Title:       source.Title,
		Description: source.Description,
		CreatedAt:   createdAt,
	}

	switch source.PublishedStatus {
	default:
	case eventProto.PublishedStatus_Draft:
		result.PublishedStatus = PublishedStatusDraft
	case eventProto.PublishedStatus_Published:
		result.PublishedStatus = PublishedStatusPublished
	}

	// result := Event{
	// 	ID:              1,
	// 	Title:           "myTitle",
	// 	Description:     "myDesc",
	// 	PublishedStatus: PublishedStatusDraft,
	// 	CreatedAt:       time.Now(),
	// }

	return &result
}
