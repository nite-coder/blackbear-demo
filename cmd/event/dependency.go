package event

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/golang-migrate/migrate/v4"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
	"github.com/jasonsoft/log/v2/handlers/gelf"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/event"
	eventGRPC "github.com/jasonsoft/starter/pkg/event/delivery/grpc"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	eventDatabase "github.com/jasonsoft/starter/pkg/event/repository/database"
	eventService "github.com/jasonsoft/starter/pkg/event/service"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	// repo
	_eventRepo event.Repository

	// services
	_eventService event.Servicer

	// grpc server
	_eventServer eventProto.EventServiceServer
)

func initialize(cfg config.Configuration) error {
	initLogger("event", cfg)

	db, err := initDatabase(cfg, "starter")
	if err != nil {
		return err
	}

	// repo
	_eventRepo = eventDatabase.NewEventRepository(cfg, db)

	// services
	_eventService = eventService.NewEventService(cfg, _eventRepo)

	// grpc server
	_eventServer = eventGRPC.NewEventServer(cfg, _eventService)

	if _eventServer == nil {
		log.Debug("event server is nil")
	}

	log.Info("event server is initialized")
	return nil
}

func initLogger(appID string, cfg config.Configuration) {
	// set up log target
	log.
		Str("app_id", appID).
		Str("env", cfg.Env).
		SaveToDefault()

	for _, target := range cfg.Logs {
		switch target.Type {
		case "console":
			clog := console.New()
			levels := log.GetLevelsFromMinLevel(target.MinLevel)
			log.AddHandler(clog, levels...)
		case "gelf":
			graylog := gelf.New(target.ConnectionString)
			levels := log.GetLevelsFromMinLevel(target.MinLevel)
			log.AddHandler(graylog, levels...)
		}
	}
}

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer(cfg config.Configuration) func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(cfg.Jaeger.AdvertiseAddr),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "event",
			Tags: []label.KeyValue{
				label.String("version", "1.0"),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Err(err).Fatal("install jaeger pipleline failed.")
	}

	return func() {
		flush()
	}
}

func grpcInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		logger := log.FromContext(ctx)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, grpc.Errorf(codes.DataLoss, "metadata is not found")
		}

		// get requestID from metadata and create a new log context
		var requestID string
		if val, ok := md["request_id"]; ok {
			requestID = val[0]
		}
		logger = logger.Str("request_id", requestID)
		ctx = logger.WithContext(ctx)

		//logger.Debugf("dump metadata %#v", md)

		// var claims identity.Claims
		// if val, ok := md["claims"]; ok {
		// 	claimsStr = val[0]

		// 	if err := json.Unmarshal([]byte(claimsStr), &claims); err != nil {

		// 	}
		// }

		//logger = log.Str("request_id", requestID).Str("claims", claims)
		//ctx = log.NewContext(ctx, logger)
		//ctx = identity.NewContext(ctx, &claims)

		// received request id
		//logger.Debugf("========== request_id: %s, claims: %s", requestID, claims)

		result, err := handler(ctx, req)
		if err != nil {
			// centralized error
			logger.Err(err).Errorf("event grpc unknown error: %v", err)
		}

		return result, err

	}
}

func initDatabase(cfg config.Configuration, name string) (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	var connectionString string
	for _, database := range cfg.Databases {
		if strings.EqualFold(database.Name, name) {

			switch strings.ToLower(database.Type) {
			case "mysql":
				connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true", database.Username, database.Password, database.Address, database.DBName)
			}

			// migrate database if needed
			if database.IsMigrated {
				path := cfg.Path("deployments", "database", database.DBName)
				path = filepath.ToSlash(path) // due to migrate package path issue on window os, therefore, we need to run this
				source := fmt.Sprintf("file://%s", path)
				migrateDBURL := fmt.Sprintf("%s://%s", database.Type, connectionString)

				m, err := migrate.New(
					source,
					migrateDBURL,
				)
				if err != nil {
					return nil, fmt.Errorf("db migration config is wrong. db: %s, source: %s, migrateDBURL: %s %w", database.DBName, source, migrateDBURL, err)
				}

				err = m.Up()
				if err != nil && !errors.Is(err, migrate.ErrNoChange) {
					return nil, fmt.Errorf("db migration failed. db: %s, source: %s, migrateDBURL: %s %w", database.DBName, source, migrateDBURL, err)
				}

				log.Infof("%s database was migrated", database.DBName)
			}
		}
	}

	var db *gorm.DB
	var err error
	err = backoff.Retry(func() error {
		db, err := gorm.Open(gormMySQL.Open(connectionString), &gorm.Config{})
		if err != nil {
			log.Errorf("main: mysql open failed: %v", err)
			return err
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
			log.Errorf("main: mysql ping error: %v", err)
			return err
		}

		return nil
	}, bo)

	if err != nil {
		log.Panicf("main: mysql connect err: %s", err.Error())
	}

	return db, nil
}
