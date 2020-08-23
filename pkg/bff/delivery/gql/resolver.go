package gql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
)

type Resolver struct {
	eventClient eventProto.EventServiceClient
}

func NewResolver(eventClient eventProto.EventServiceClient) *Resolver {
	return &Resolver{
		eventClient: eventClient,
	}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
