package gql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
	temporalClient "go.temporal.io/sdk/client"
)

type Resolver struct {
	eventClient    eventProto.EventServiceClient
	walletClient   walletProto.WalletServiceClient
	temporalClient temporalClient.Client
}

func NewResolver(eventClient eventProto.EventServiceClient, walletClient walletProto.WalletServiceClient, temporalClient temporalClient.Client) *Resolver {
	return &Resolver{
		eventClient:    eventClient,
		walletClient:   walletClient,
		temporalClient: temporalClient,
	}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
