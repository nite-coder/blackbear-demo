package gql

import (
	"context"

	"github.com/nite-coder/blackbear-demo/pkg/event/proto"
	"github.com/nite-coder/blackbear/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *queryResolver) GetEvents(ctx context.Context, input GetEventOptionsInput) ([]*Event, error) {
	logger := log.FromContext(ctx)
	logger.Debug("gql: begin get events fn")

	resp, err := r.eventClient.GetEvents(ctx, &proto.GetEventsRequest{
		Id:    input.ID,
		Title: input.Title,
	})
	if err != nil {
		return nil, err
	}

	result := []*Event{}
	for _, data := range resp.Events {
		event, err := eventToGQL(data)
		if err != nil {
			return nil, err
		}
		result = append(result, event)
	}

	return result, nil
}

func (r *queryResolver) GetEvent(ctx context.Context, eventID int64) (*Event, error) {
	logger := log.FromContext(ctx)
	logger.Debug("gql: begin get event fn")

	resp, err := r.eventClient.GetEvents(ctx, &proto.GetEventsRequest{
		Id: eventID,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Events) == 0 {
		return nil, ErrNotFound
	}

	event, err := eventToGQL(resp.Events[0])
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *queryResolver) GetWallet(ctx context.Context) (*Wallet, error) {
	logger := log.FromContext(ctx)
	logger.Debug("gql: begin get wallet fn")

	resp, err := r.walletClient.GetWallet(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	result, err := walletToGQL(resp.Data)
	if err != nil {
		return nil, err
	}

	return result, nil

}
