package event

import (
	"context"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/event"
	eventGRPC "github.com/jasonsoft/starter/pkg/event/delivery/grpc"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	eventDatabase "github.com/jasonsoft/starter/pkg/event/repository/database"
	eventService "github.com/jasonsoft/starter/pkg/event/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

	db, err := cfg.InitDatabase("starter_db")
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
