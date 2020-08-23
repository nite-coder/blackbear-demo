package event

import (
	"context"
	"time"
)

type PublishedStatus int32

const (
	Draft     PublishedStatus = 0
	Published PublishedStatus = 1
)

type Event struct {
	ID              int64
	Title           string
	Description     string
	PublishedStatus PublishedStatus
	CreatedAt       time.Time
}

type UpdateEventStatusRequest struct {
	EventID         int64
	TransID         string
	PublishedStatus PublishedStatus
}

type EventServicer interface {
	Events(ctx context.Context) ([]*Event, error)
	UpdatePublishStatus(ctx context.Context, request UpdateEventStatusRequest) error
}
