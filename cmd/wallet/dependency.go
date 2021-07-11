package wallet

import (
	"github.com/nite-coder/blackbear-demo/internal/pkg/config"
	"github.com/nite-coder/blackbear-demo/pkg/domain"
	walletGRPC "github.com/nite-coder/blackbear-demo/pkg/wallet/delivery/grpc"
	walletProto "github.com/nite-coder/blackbear-demo/pkg/wallet/proto"
	walletUsecase "github.com/nite-coder/blackbear-demo/pkg/wallet/usecase"
	"github.com/nite-coder/blackbear/pkg/log"
)

var (
	_walletService domain.WalletUsecase

	// grpc server
	_walletServer walletProto.WalletServiceServer
)

func initialize(cfg config.Configuration) error {
	cfg.InitLogger("wallet")

	_walletService = walletUsecase.NewWalletUsecase(cfg)

	_walletServer = walletGRPC.NewWalletServer(cfg, _walletService)

	if _walletServer == nil {
		log.Debug("wallet server is nil")
	}

	log.Info("wallet server is initialized")
	return nil
}
