package frontend

import (
	"github.com/nite-coder/blackbear-demo/internal/pkg/config"
	eventGRPC "github.com/nite-coder/blackbear-demo/pkg/event/delivery/grpc"
	eventProto "github.com/nite-coder/blackbear-demo/pkg/event/proto"
	walletGRPC "github.com/nite-coder/blackbear-demo/pkg/wallet/delivery/grpc"
	walletProto "github.com/nite-coder/blackbear-demo/pkg/wallet/proto"
	starterWorkflow "github.com/nite-coder/blackbear-demo/pkg/workflow"
	"github.com/nite-coder/blackbear/pkg/log"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/bridge/opentracing"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/client"
	temporalClient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
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

	_tracer = otel.Tracer("")

	_eventClient, err = eventGRPC.NewClient(cfg)
	if err != nil {
		return err
	}

	_walletClient, err = walletGRPC.NewClient(cfg)
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
