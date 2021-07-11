package gql

import (
	"github.com/golang/protobuf/ptypes"
	eventProto "github.com/nite-coder/blackbear-demo/pkg/event/proto"
	walletProto "github.com/nite-coder/blackbear-demo/pkg/wallet/proto"
)

func eventToGQL(event *eventProto.Event) (*Event, error) {
	if event == nil {
		return nil, nil
	}

	createdAt, err := ptypes.Timestamp(event.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := ptypes.Timestamp(event.UpdatedAt)
	if err != nil {
		return nil, err
	}

	result := Event{
		ID:          event.Id,
		Title:       event.Title,
		Description: event.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	switch event.PublishedStatus {
	default:
	case eventProto.PublishedStatus_PublishedStatus_Draft:
		result.PublishedStatus = PublishedStatusDraft
	case eventProto.PublishedStatus_PublishedStatus_Published:
		result.PublishedStatus = PublishedStatusPublished
	}

	return &result, nil
}

func walletToGQL(wallet *walletProto.Wallet) (*Wallet, error) {
	if wallet == nil {
		return nil, nil
	}

	updatedAt, err := ptypes.Timestamp(wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}

	result := Wallet{
		ID:        wallet.Id,
		Amount:    wallet.Amount,
		UpdatedAt: updatedAt,
	}

	return &result, nil
}
