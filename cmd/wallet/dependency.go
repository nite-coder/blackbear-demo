package wallet

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/wallet"
	walletGRPC "github.com/jasonsoft/starter/pkg/wallet/delivery/grpc"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
	walletService "github.com/jasonsoft/starter/pkg/wallet/service"
)

var (
	_walletService wallet.Servicer

	// grpc server
	_walletServer walletProto.WalletServiceServer
)

func initialize(cfg config.Configuration) error {
	cfg.InitLogger("wallet")

	_walletService = walletService.NewWalletService(cfg)

	_walletServer = walletGRPC.NewWalletServer(cfg, _walletService)

	if _walletServer == nil {
		log.Debug("wallet server is nil")
	}

	log.Info("wallet server is initialized")
	return nil
}
