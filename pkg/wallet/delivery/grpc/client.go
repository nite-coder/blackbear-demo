package grpc

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	internalGRPC "github.com/jasonsoft/starter/internal/pkg/grpc"
	"github.com/jasonsoft/starter/pkg/wallet/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func NewClient(cfg config.Configuration) (proto.WalletServiceClient, error) {

	conn, err := grpc.Dial(cfg.Wallet.GRPCAdvertiseAddr,
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
		log.Err(err).Errorf("grpc: dial to wallet server failed. connection string: %s", cfg.Wallet.GRPCAdvertiseAddr)
		return nil, err
	}

	log.Info("grpc: dail to wallet server connect successfully")

	client := proto.NewWalletServiceClient(conn)
	return client, nil
}
