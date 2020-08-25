package gql

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jasonsoft/log/v2"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *queryResolver) GetEvents(ctx context.Context) ([]*Event, error) {
	logger := log.FromContext(ctx)
	logger.Debug("gql: begin get events fn")

	resp, err := r.eventClient.GetEvents(ctx, &emptypb.Empty{})
	if err != nil {
		err = fmt.Errorf("eventClient call failed, name: %s, %w", "GetEvents", err)
		logger.Err(err).Warn("gql: eventClient call failed")
		return nil, err
	}

	result := []*Event{}

	for _, d := range resp.Data {
		result = append(result, grpcEventToGQLEvent(d))
	}

	return result, nil
}

func (r *queryResolver) GetWallet(ctx context.Context) (*Wallet, error) {
	logger := log.FromContext(ctx)
	logger.Debug("gql: begin get wallet fn")

	resp, err := r.walletClient.GetWallet(ctx, &emptypb.Empty{})
	if err != nil {
		err = fmt.Errorf("walletClient call failed, name: %s, %w", "GetWallet", err)
		logger.Err(err).Warn("gql: walletClient call failed")
		return nil, err
	}

	result := grpcWalletToGQLWallet(resp.Data)

	return result, nil

}

func grpcEventToGQLEvent(source *eventProto.Event) *Event {
	if source == nil {
		return nil
	}

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

	return &result
}

func grpcWalletToGQLWallet(source *walletProto.Wallet) *Wallet {
	if source == nil {
		return nil
	}

	updatedAt, _ := ptypes.Timestamp(source.UpdatedAt)

	result := Wallet{
		ID:        source.Id,
		Amount:    source.Amount,
		UpdatedAt: updatedAt,
	}

	return &result
}
