package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nite-coder/blackbear-demo/internal/pkg/config"
	"github.com/nite-coder/blackbear-demo/pkg/domain"
	"gorm.io/gorm"
)

// EventRepo handles all event database operations
type EventRepo struct {
	config config.Configuration
	db     *gorm.DB
}

func NewEventRepository(cfg config.Configuration, db *gorm.DB) domain.EventRepository {
	return &EventRepo{
		config: cfg,
		db:     db,
	}
}

func (repo *EventRepo) WithTX(tx *gorm.DB) domain.EventRepository {
	if tx == nil {
		return repo
	}
	newRepo := *repo
	newRepo.db = tx
	return &newRepo
}

func (repo *EventRepo) Event(ctx context.Context, id int64) (domain.Event, error) {
	db := repo.db.WithContext(ctx)

	event := domain.Event{}

	err := db.First(&event, "id = @id", sql.Named("id", id)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return event, domain.ErrNotFound
		}
		return event, err
	}

	return event, nil
}

func (repo *EventRepo) Events(ctx context.Context, opts domain.FindEventOptions) ([]domain.Event, error) {
	db := repo.db.WithContext(ctx)

	events := []domain.Event{}

	db = repo.buildSQL(ctx, db, opts)
	err := db.Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (repo *EventRepo) buildSQL(ctx context.Context, db *gorm.DB, opts domain.FindEventOptions) *gorm.DB {

	// if opts.ID > 0 {
	// 	db = db.Where("id = @id", sql.Named("id", opts.ID))
	// }

	// if len(opts.Title) > 0 {
	// 	db = db.Where("title = @title", sql.Named("title", opts.Title))
	// }
	db = db.Where(opts)

	if opts.CreatedAtStart.IsZero() == false && opts.CreatedAtEnd.IsZero() == false {
		db = db.
			Where("created_at >= @created_at", sql.Named("created_at_start", opts.CreatedAtStart)).
			Where("created_at < @created_at_end", sql.Named("created_at_end", opts.CreatedAtEnd))
	}

	return db

}

func (repo *EventRepo) UpdatePublishStatus(ctx context.Context, request domain.UpdateEventStatusRequest) error {
	db := repo.db.WithContext(ctx)

	result := db.Model(domain.Event{}).
		Where("id = @id", sql.Named("id", request.EventID)).
		Where("version = @version", sql.Named("version", request.Version)).
		UpdateColumn("published_status", request.PublishedStatus).
		UpdateColumn("version", gorm.Expr("version + 1"))

	if result.RowsAffected == 0 {
		return domain.ErrStale
	}

	return result.Error
}
