package domain

import (
	"context"
	"time"

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

// Event is Event
type Event struct {
	ID              int64           `gorm:"column:id;primary_key" json:"id"`
	Title           string          `gorm:"column:title" json:"title"`
	Description     string          `gorm:"column:description" json:"description"`
	PublishedStatus PublishedStatus `gorm:"column:published_status" json:"published_status"`
	Version         int64           `gorm:"column:version" json:"version"`
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
	Version         int64
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
	WithTX(tx *gorm.DB) EventRepository
	Event(ctx context.Context, id int64) (Event, error)
	Events(ctx context.Context, opts FindEventOptions) ([]Event, error)
	UpdatePublishStatus(ctx context.Context, request UpdateEventStatusRequest) error
}
