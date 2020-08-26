package event

import (
	"context"
	"time"
)

// PublishedStatus identifies event published status.
type PublishedStatus int32

const (
	// Draft ...
	Draft PublishedStatus = 0
	// Published ...
	Published PublishedStatus = 1
)

// Event is Event
type Event struct {
	ID              int64
	Title           string
	Description     string
	PublishedStatus PublishedStatus
	CreatedAt       time.Time
}

// UpdateEventStatusRequest ...
type UpdateEventStatusRequest struct {
	EventID         int64
	TransID         string
	PublishedStatus PublishedStatus
}

// Servicer handles event's business logic
type Servicer interface {
	Events(ctx context.Context) ([]*Event, error)
	UpdatePublishStatus(ctx context.Context, request UpdateEventStatusRequest) error
}
