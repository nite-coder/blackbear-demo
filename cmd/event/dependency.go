package event

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/internal/pkg/database"
	"github.com/jasonsoft/starter/pkg/domain"
	eventGRPC "github.com/jasonsoft/starter/pkg/event/delivery/grpc"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	eventDatabase "github.com/jasonsoft/starter/pkg/event/repository/mysql"
	eventUsecase "github.com/jasonsoft/starter/pkg/event/usecase"
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
