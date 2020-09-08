package domain

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

// PublishedStatus identifies event published status.
type PublishedStatus int32

const (
	// Draft ...
	Draft PublishedStatus = 1
	// Published ...
	Published PublishedStatus = 2
)

var (
	// ErrNotFound means resource not found
	ErrNotFound = &AppError{Code: "NOT_FOUND", Message: "resource was not found or status was wrong", Status: codes.NotFound}
)

// Event is Event
type Event struct {
	ID              int64           `gorm:"column:id;primary_key" json:"id"`
	Title           string          `gorm:"column:title" json:"title"`
	Description     string          `gorm:"column:description" json:"description"`
	PublishedStatus PublishedStatus `gorm:"column:published_status" json:"published_status"`
	CreatedAt       time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (e *Event) TableName() string {
	return "events"
}

// UpdateEventStatusRequest ...
type UpdateEventStatusRequest struct {
	EventID         int64
	TransID         string
	PublishedStatus PublishedStatus
}

// FindEventOptions is a query condition for finding events
type FindEventOptions struct {
	ID             int64  `gorm:"column:id;primary_key" json:"id"`
	Title          string `gorm:"column:title" json:"title"`
	CreatedAtStart time.Time
	CreatedAtEnd   time.Time
}

func (f *FindEventOptions) TableName() string {
	return "events"
}

// EventUsecase represents the wallet's usecases
type EventUsecase interface {
	Events(ctx context.Context, opts FindEventOptions) ([]Event, error)
	UpdatePublishStatus(ctx context.Context, request UpdateEventStatusRequest) error
}

// EventRepository represents the event's repository contract
type EventRepository interface {
	Events(ctx context.Context, opts FindEventOptions, tx ...*gorm.DB) ([]Event, error)
	UpdatePublishStatus(ctx context.Context, request UpdateEventStatusRequest, tx ...*gorm.DB) error
}
