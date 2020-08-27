package database

import (
	"context"
	"errors"

	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/event"
	"gorm.io/gorm"
)

// EventRepo handles all event database operations
type EventRepo struct {
	config config.Configuration
	db     *gorm.DB
}

func NewEventRepository(cfg config.Configuration, db *gorm.DB) event.Repository {
	return &EventRepo{
		config: cfg,
		db:     db,
	}
}

func (repo *EventRepo) Events(ctx context.Context, opts event.FindEventOptions) ([]event.Event, error) {
	events := []event.Event{}

	db := repo.buildSQL(ctx, repo.db, opts)

	err := db.Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (repo *EventRepo) buildSQL(ctx context.Context, db *gorm.DB, opts event.FindEventOptions) *gorm.DB {

	if opts.ID > 0 {
		db = db.Where("id = ?", opts.ID)
	}

	if len(opts.Title) > 0 {
		db = db.Where("title = ?", opts.Title)
	}

	return db

}

func (repo *EventRepo) UpdatePublishStatus(ctx context.Context, request event.UpdateEventStatusRequest) error {

	err := repo.db.Model(event.Event{}).
		UpdateColumn("published_status", request.PublishedStatus).
		Where("id = ?", request.EventID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return event.ErrNotFound
		}
		return err
	}

	return nil
}
