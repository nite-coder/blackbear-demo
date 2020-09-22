package database

import (
	"context"

	"gorm.io/gorm"
)

func ExecuteTx(ctx context.Context, db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.Begin()
	return executeInTx(ctx, tx, func() error { return fn(tx) })
}

func executeInTx(ctx context.Context, tx *gorm.DB, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	return fn()
}
