package service

import (
	"context"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/domain"
)

// EventService handles event's business logic
type EventService struct {
	config config.Configuration
	repo   domain.EventRepository
}

// NewEventService create an instance of event service
func NewEventService(cfg config.Configuration, repo domain.EventRepository) domain.EventServicer {
	return &EventService{
		config: cfg,
		repo:   repo,
	}
}

// Events returns all events
func (srv *EventService) Events(ctx context.Context, opts domain.FindEventOptions) ([]domain.Event, error) {
	return srv.repo.Events(ctx, opts)
}

// UpdatePublishStatus update
func (srv *EventService) UpdatePublishStatus(ctx context.Context, request domain.UpdateEventStatusRequest) error {
	logger := log.FromContext(ctx)
	logger.Debug("service: begin UpdatePublishStatus fn")

	return srv.repo.UpdatePublishStatus(ctx, request)
}
