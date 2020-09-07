package event

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/internal/pkg/database"
	"github.com/jasonsoft/starter/pkg/event"
	eventGRPC "github.com/jasonsoft/starter/pkg/event/delivery/grpc"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	eventDatabase "github.com/jasonsoft/starter/pkg/event/repository/database"
	eventService "github.com/jasonsoft/starter/pkg/event/service"
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
	cfg.InitLogger("event")

	db, err := database.InitDatabase(cfg, "starter_db")
	if err != nil {
		return err
	}

	// repo
	_eventRepo = eventDatabase.NewEventRepository(cfg, db)

	// services
	_eventService = eventService.NewEventService(cfg, _eventRepo)

	// grpc server
	_eventServer = eventGRPC.NewEventServer(cfg, _eventService)

	log.Info("event server is initialized")
	return nil
}
