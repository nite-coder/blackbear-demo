package usecase

import (
	"context"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/internal/pkg/database"
	"github.com/jasonsoft/starter/pkg/domain"
	"gorm.io/gorm"
)

type eventUsecase struct {
	config config.Configuration
	db     *gorm.DB
	repo   domain.EventRepository
}

// NewEventUsecase create a new eventUsecase object representation of domain.eventUsecase interface
func NewEventUsecase(cfg config.Configuration, db *gorm.DB, repo domain.EventRepository) domain.EventUsecase {
	return &eventUsecase{
		config: cfg,
		db:     db,
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

	err := database.ExecuteTx(ctx, u.db, func(tx *gorm.DB) error {
		eventRepo := u.repo.WithTX(tx)

		event, err := eventRepo.Event(ctx, request.EventID)
		if err != nil {
			return err
		}

		if event.PublishedStatus == domain.Published {
			return domain.ErrWrongStatus
		}

		request.Version = event.Version
		return eventRepo.UpdatePublishStatus(ctx, request)
	})

	if err != nil {
		return err
	}

	return nil
}
