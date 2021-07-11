package event

import (
	"github.com/nite-coder/blackbear-demo/internal/pkg/config"
	"github.com/nite-coder/blackbear-demo/internal/pkg/database"
	"github.com/nite-coder/blackbear-demo/pkg/domain"
	eventGRPC "github.com/nite-coder/blackbear-demo/pkg/event/delivery/grpc"
	eventProto "github.com/nite-coder/blackbear-demo/pkg/event/proto"
	eventDatabase "github.com/nite-coder/blackbear-demo/pkg/event/repository/mysql"
	eventUsecase "github.com/nite-coder/blackbear-demo/pkg/event/usecase"
	"github.com/nite-coder/blackbear/pkg/log"
)

var (
	// repo
	_eventRepo domain.EventRepository

	// services
	_eventService domain.EventUsecase

	// grpc server
	_eventServer eventProto.EventServiceServer
)

func initialize(cfg config.Configuration) error {
	cfg.InitLogger("event")

	db, err := database.InitDatabase(cfg, "starter_db")
	if err != nil {
		return err
	}

	// repo
	_eventRepo = eventDatabase.NewEventRepository(cfg, db)

	// services
	_eventService = eventUsecase.NewEventUsecase(cfg, db, _eventRepo)

	// grpc server
	_eventServer = eventGRPC.NewEventServer(cfg, _eventService)

	log.Info("event server is initialized")
	return nil
}
