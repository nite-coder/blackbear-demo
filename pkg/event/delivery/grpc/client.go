package grpc

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	internalGRPC "github.com/jasonsoft/starter/internal/pkg/grpc"
	"github.com/jasonsoft/starter/pkg/event/proto"
	grpctrace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func NewClient(cfg config.Configuration) (proto.EventServiceClient, error) {
	tracer := global.Tracer("")

	conn, err := grpc.Dial(cfg.Event.GRPCAdvertiseAddr,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5,
			Timeout:             5,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(
			grpctrace.UnaryClientInterceptor(tracer),
			internalGRPC.ClientInterceptor(),
		),
		grpc.WithStreamInterceptor(grpctrace.StreamClientInterceptor(tracer)),
	)

	if err != nil {
		log.Err(err).Errorf("grpc: dial to event server failed. connection string: %s", cfg.Event.GRPCAdvertiseAddr)
		return nil, err
	}

	log.Info("grpc: dail to event server connect successfully")

	client := proto.NewEventServiceClient(conn)
	return client, nil
}
