package event

import (
	"context"
	"time"

	"github.com/jasonsoft/starter/internal/pkg/exception"
)

// PublishedStatus identifies event published status.
type PublishedStatus int32

const (
	// Draft ...
	Draft PublishedStatus = 0
	// Published ...
	Published PublishedStatus = 1
)

var (
	// ErrNotFound means resource not found
	ErrNotFound = exception.New("NOT_FOUND", "event resource was not found")
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
	ID              int64
	Title           string
	PublishedStatus PublishedStatus
	CreatedAt       time.Time
}

// Servicer handles event's business logic
type Servicer interface {
	Events(ctx context.Context, opts FindEventOptions) ([]Event, error)
	UpdatePublishStatus(ctx context.Context, request UpdateEventStatusRequest) error
}

// Repository handles event's database operations
type Repository interface {
	Events(ctx context.Context, opts FindEventOptions) ([]Event, error)
	UpdatePublishStatus(ctx context.Context, request UpdateEventStatusRequest) error
}
