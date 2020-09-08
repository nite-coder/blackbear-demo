package wallet

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/domain"
	walletGRPC "github.com/jasonsoft/starter/pkg/wallet/delivery/grpc"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
	walletUsecase "github.com/jasonsoft/starter/pkg/wallet/usecase"
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
