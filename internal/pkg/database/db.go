package database

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(cfg config.Configuration, name string) (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	for _, database := range cfg.Databases {
		if strings.EqualFold(database.Name, name) {

			// migrate database if needed
			if database.IsMigrated {
				path := cfg.Path("deployments", "database", database.Name)
				path = filepath.ToSlash(path) // due to migrate package path issue on window os, therefore, we need to run this
				source := fmt.Sprintf("file://%s", path)
				migrateDBURL := fmt.Sprintf("%s://%s", database.Type, database.ConnectionString)

				m, err := migrate.New(
					source,
					migrateDBURL,
				)
				if err != nil {
					// make sure migration package is using v4 above
					return nil, fmt.Errorf("db migration config was wrong. db_name: %s, source: %s, migrateDBURL: %s, error: %w", database.Name, source, migrateDBURL, err)
				}

				err = m.Up()
				if err != nil && !errors.Is(err, migrate.ErrNoChange) {
					return nil, fmt.Errorf("db migration failed. db: %s, source: %s, migrateDBURL: %s, error: %w", database.Name, source, migrateDBURL, err)
				}

				log.Infof("%s database was migrated", database.Name)
			}

			var db *gorm.DB
			var err error
			err = backoff.Retry(func() error {

				gormConfig := gorm.Config{
					//PrepareStmt: true,
					Logger: logger.Default.LogMode(logger.Silent),
				}

				switch strings.ToLower(database.Type) {
				case "mysql":
					db, err = gorm.Open(gormMySQL.Open(database.ConnectionString), &gormConfig)
				}

				if err != nil {
					return fmt.Errorf("database: database open failed: %w", err)
				}

				err = db.Use(&TelemetryPlugin{})
				if err != nil {
					log.Err(err).Warn("database: register telemetry plugin failed")
				}

				sqlDB, err := db.DB()
				if err != nil {
					return err
				}

				sqlDB.SetMaxIdleConns(150)
				sqlDB.SetMaxOpenConns(300)
				sqlDB.SetConnMaxLifetime(14400 * time.Second)

				err = sqlDB.Ping()
				if err != nil {
					return fmt.Errorf("database: database ping failed. name: %s, error: %w", name, err)
				}

				return nil
			}, bo)

			if err != nil {
				return nil, fmt.Errorf("database: database connect failed.  name: %s, error: %w", name, err)
			}

			return db, nil
		}
	}

	return nil, fmt.Errorf("database: database name was not found. name: %s", name)

}

func RunSQLScripts(db *gorm.DB, dirPath string) error {

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		// TODO: it seems to me that if we turn gorm prepare config to true, we will receive the error here.  Now, we just turn that off.
		err = db.Exec(string(data)).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
