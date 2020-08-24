package bff

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
	"github.com/jasonsoft/log/v2/handlers/gelf"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/napnap/middleware"
	"github.com/jasonsoft/starter/internal/pkg/config"
	internalMiddleware "github.com/jasonsoft/starter/internal/pkg/middleware"
	"github.com/jasonsoft/starter/pkg/bff/delivery/gql"
	bffGRPC "github.com/jasonsoft/starter/pkg/bff/delivery/grpc"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/instrumentation/grpctrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	_tracer trace.Tracer
	// grpc clients
	_walletClient walletProto.WalletServiceClient
	_eventClient  eventProto.EventServiceClient
)

func initialize(cfg config.Configuration) error {
	var err error

	initLogger("bff", cfg)

	_tracer = global.Tracer("")

	_eventClient, err = eventGRPCClient(cfg)
	if err != nil {
		return err
	}

	_walletClient, err = walletGRPCClient(cfg)
	if err != nil {
		return err
	}

	log.Info("bff server is initialized")
	return nil
}

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer(cfg config.Configuration) func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(cfg.Jaeger.AdvertiseAddr),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "bff",
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

func eventGRPCClient(cfg config.Configuration) (eventProto.EventServiceClient, error) {
	conn, err := grpc.Dial(cfg.Event.GRPCAdvertiseAddr,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5,
			Timeout:             5,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(
			grpctrace.UnaryClientInterceptor(_tracer),
			bffGRPC.ClientInterceptor(),
		),
		grpc.WithStreamInterceptor(grpctrace.StreamClientInterceptor(_tracer)),
	)

	if err != nil {
		log.Errorf("main: dial event grpc server failed: %v, connection string: %s", err, cfg.Event.GRPCAdvertiseAddr)
		return nil, err
	}

	log.Infof("main: dail event grpc server %s%s", cfg.Event.GRPCAdvertiseAddr, " connect successfully")

	client := eventProto.NewEventServiceClient(conn)
	return client, nil
}

func walletGRPCClient(cfg config.Configuration) (walletProto.WalletServiceClient, error) {
	conn, err := grpc.Dial(cfg.Wallet.GRPCAdvertiseAddr,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5,
			Timeout:             5,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(
			grpctrace.UnaryClientInterceptor(_tracer),
			bffGRPC.ClientInterceptor(),
		),
		grpc.WithStreamInterceptor(grpctrace.StreamClientInterceptor(_tracer)),
	)

	if err != nil {
		log.Errorf("main: dial wallet grpc server failed: %v, connection string: %s", err, cfg.Wallet.GRPCAdvertiseAddr)
		return nil, err
	}

	log.Infof("main: dail wallet grpc server %s%s", cfg.Wallet.GRPCAdvertiseAddr, " connect successfully")

	client := walletProto.NewWalletServiceClient(conn)
	return client, nil
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

func verifyOrigin(origin string) bool {
	return true
}

func newNapNap() *napnap.NapNap {
	nap := napnap.New()
	nap.Use(internalMiddleware.NewRequestIDMW())
	nap.Use(internalMiddleware.NewTracerMW())
	nap.Use(internalMiddleware.NewLoggerMW())

	// turn on CORS feature
	options := middleware.Options{}
	options.AllowOriginFunc = verifyOrigin
	options.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTION", "HEAD"}
	options.AllowedHeaders = []string{"*", "Authorization", "Content-Type", "Origin", "Content-Length", "accept"}
	nap.Use(middleware.NewCors(options))

	nap.Get("/", func(c *napnap.Context) error {
		return c.String(200, "Hello World")
	})

	nap.Get("/ping", func(c *napnap.Context) error {
		return c.String(http.StatusOK, "pong!!!")
	})

	nap.Get("/playground", napnap.WrapHandler(handler.Playground("GraphQL playground", "/graphql")))

	rootResolver := gql.NewResolver(_eventClient, _walletClient)

	nap.Post("/graphql", napnap.WrapHandler(handler.GraphQL(
		gql.NewExecutableSchema(gql.Config{Resolvers: rootResolver}),
		handler.ErrorPresenter(
			func(ctx context.Context, e error) *gqlerror.Error {
				logger := log.FromContext(ctx)

				// TODO: handle all graphql errors here
				err := &gqlerror.Error{
					Message: e.Error(),
					Extensions: map[string]interface{}{
						"code": "401001",
					},
				}

				logger.Err(e).Error("GQL unknown error")
				return err
			}),
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			logger := log.FromContext(ctx)
			myErr := fmt.Errorf("Internal server error! %w", err)

			logger.Err(myErr).Error("Internal server error!")
			return myErr
		}),
	)))

	return nap
}
