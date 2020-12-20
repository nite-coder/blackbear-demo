package grpc

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	internalGRPC "github.com/jasonsoft/starter/internal/pkg/grpc"
	"github.com/jasonsoft/starter/pkg/event/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func NewClient(cfg config.Configuration) (proto.EventServiceClient, error) {

	conn, err := grpc.Dial(cfg.Event.GRPCAdvertiseAddr,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5,
			Timeout:             5,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
			internalGRPC.ClientInterceptor(),
		),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)

	if err != nil {
		log.Err(err).Errorf("grpc: dial to event server failed. connection string: %s", cfg.Event.GRPCAdvertiseAddr)
		return nil, err
	}

	log.Info("grpc: dail to event server connect successfully")

	client := proto.NewEventServiceClient(conn)
	return client, nil
}
