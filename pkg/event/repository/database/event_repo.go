package database

import (
	"context"
	"database/sql"

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

func (repo *EventRepo) Events(ctx context.Context, opts event.FindEventOptions, tx ...*gorm.DB) ([]event.Event, error) {
	db := repo.db
	if tx != nil {
		db = tx[0]
	}
	db = db.WithContext(ctx)

	events := []event.Event{}

	db = repo.buildSQL(ctx, db, opts)
	err := db.Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (repo *EventRepo) buildSQL(ctx context.Context, db *gorm.DB, opts event.FindEventOptions) *gorm.DB {

	if opts.ID > 0 {
		db = db.Where("id = @id", sql.Named("id", opts.ID))
	}

	if len(opts.Title) > 0 {
		db = db.Where("title = @title", sql.Named("title", opts.Title))
	}

	return db

}

func (repo *EventRepo) UpdatePublishStatus(ctx context.Context, request event.UpdateEventStatusRequest, tx ...*gorm.DB) error {
	db := repo.db
	if tx != nil {
		db = tx[0]
	}
	db = db.WithContext(ctx)

	result := db.Model(event.Event{}).
		Where("id = @id", sql.Named("id", request.EventID)).
		Where("published_status = @published_status", sql.Named("published_status", event.Draft)).
		UpdateColumn("published_status", request.PublishedStatus)

	if result.RowsAffected == 0 {
		return event.ErrNotFound
	}

	return result.Error
}
