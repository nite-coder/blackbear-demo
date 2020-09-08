package frontend

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	frontendGRPC "github.com/jasonsoft/starter/pkg/frontend/delivery/grpc"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
	starterWorkflow "github.com/jasonsoft/starter/pkg/workflow"
	grpctrace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/bridge/opentracing"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.temporal.io/sdk/client"
	temporalClient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	_tracer         trace.Tracer
	_temporalClient temporalClient.Client
	// grpc clients
	_walletClient walletProto.WalletServiceClient
	_eventClient  eventProto.EventServiceClient
)

func initialize(cfg config.Configuration) error {
	var err error

	cfg.InitLogger("frontend")

	_tracer = global.Tracer("")

	_eventClient, err = eventGRPCClient(cfg)
	if err != nil {
		return err
	}

	_walletClient, err = walletGRPCClient(cfg)
	if err != nil {
		return err
	}

	_temporalClient, err = initTemporalClient(cfg)
	if err != nil {
		return err
	}

	log.Info("frontend server is initialized")
	return nil
}

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer(cfg config.Configuration) func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(cfg.Jaeger.AdvertiseAddr),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "frontend",
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

func initTemporalClient(cfg config.Configuration) (temporalClient.Client, error) {
	// The client is a heavyweight object that should be created once per process.
	bridgeTracer, _ := opentracing.NewTracerPair(_tracer)

	c, err := temporalClient.NewClient(client.Options{
		HostPort: cfg.Temporal.Address,
		ContextPropagators: []workflow.ContextPropagator{
			starterWorkflow.NewContextPropagator(),
		},
		Tracer: bridgeTracer,
	})
	if err != nil {
		return nil, err
	}
	return c, nil
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
			frontendGRPC.ClientInterceptor(),
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
			frontendGRPC.ClientInterceptor(),
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
