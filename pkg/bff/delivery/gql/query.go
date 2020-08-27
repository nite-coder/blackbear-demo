package gql

import (
	"context"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/pkg/event/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *queryResolver) GetEvents(ctx context.Context, input *GetEventOptionsInput) ([]*Event, error) {
	logger := log.FromContext(ctx)
	logger.Debug("gql: begin get events fn")

	resp, err := r.eventClient.GetEvents(ctx, &proto.GetEventsRequest{})
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

func (r *queryResolver) GetEvent(ctx context.Context, eventID *int64) (*Event, error) {
	panic("not implemented")
}
