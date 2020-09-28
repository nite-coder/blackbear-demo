package grpc

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	internalGRPC "github.com/jasonsoft/starter/internal/pkg/grpc"
	"github.com/jasonsoft/starter/pkg/wallet/proto"
	grpctrace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/api/global"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func NewClient(cfg config.Configuration) (proto.WalletServiceClient, error) {
	tracer := global.Tracer("")

	conn, err := grpc.Dial(cfg.Wallet.GRPCAdvertiseAddr,
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
		log.Err(err).Errorf("grpc: dial to wallet server failed. connection string: %s", cfg.Wallet.GRPCAdvertiseAddr)
		return nil, err
	}

	log.Info("grpc: dail to wallet server connect successfully")

	client := proto.NewWalletServiceClient(conn)
	return client, nil
}
