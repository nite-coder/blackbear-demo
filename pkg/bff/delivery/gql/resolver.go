package gql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
)

type Resolver struct {
	eventClient  eventProto.EventServiceClient
	walletClient walletProto.WalletServiceClient
}

func NewResolver(eventClient eventProto.EventServiceClient, walletClient walletProto.WalletServiceClient) *Resolver {
	return &Resolver{
		eventClient:  eventClient,
		walletClient: walletClient,
	}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
