package usecase

import (
	"context"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/domain"
)

type eventUsecase struct {
	config config.Configuration
	repo   domain.EventRepository
}

// NewEventUsecase create a new eventUsecase object representation of domain.eventUsecase interface
func NewEventUsecase(cfg config.Configuration, repo domain.EventRepository) domain.EventUsecase {
	return &eventUsecase{
		config: cfg,
		repo:   repo,
	}
}

// Events returns all events
func (u *eventUsecase) Events(ctx context.Context, opts domain.FindEventOptions) ([]domain.Event, error) {
	return u.repo.Events(ctx, opts)
}

// UpdatePublishStatus update
func (u *eventUsecase) UpdatePublishStatus(ctx context.Context, request domain.UpdateEventStatusRequest) error {
	logger := log.FromContext(ctx)
	logger.Debug("service: begin UpdatePublishStatus fn")

	return u.repo.UpdatePublishStatus(ctx, request)
}
