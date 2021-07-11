package gql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	eventProto "github.com/nite-coder/blackbear-demo/pkg/event/proto"
	walletProto "github.com/nite-coder/blackbear-demo/pkg/wallet/proto"
	temporalClient "go.temporal.io/sdk/client"
)

// Resolver will be used to handles all business logic
type Resolver struct {
	eventClient    eventProto.EventServiceClient
	walletClient   walletProto.WalletServiceClient
	temporalClient temporalClient.Client
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
