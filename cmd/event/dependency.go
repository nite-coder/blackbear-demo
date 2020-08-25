package event

import (
	"context"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
	"github.com/jasonsoft/log/v2/handlers/gelf"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/event"
	eventGRPC "github.com/jasonsoft/starter/pkg/event/delivery/grpc"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	eventService "github.com/jasonsoft/starter/pkg/event/service"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	_eventService event.EventServicer

	// grpc server
	_eventServer eventProto.EventServiceServer
)

func initialize(cfg config.Configuration) error {
	initLogger("event", cfg)

	_eventService = eventService.NewEventService(cfg)

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
			Tags: []kv.KeyValue{
				kv.String("version", "1.0"),
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
