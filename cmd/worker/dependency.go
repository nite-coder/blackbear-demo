package worker

import (
	"github.com/nite-coder/blackbear-demo/internal/pkg/config"
	eventGRPC "github.com/nite-coder/blackbear-demo/pkg/event/delivery/grpc"
	eventProto "github.com/nite-coder/blackbear-demo/pkg/event/proto"
	walletGRPC "github.com/nite-coder/blackbear-demo/pkg/wallet/delivery/grpc"
	walletProto "github.com/nite-coder/blackbear-demo/pkg/wallet/proto"
	"github.com/nite-coder/blackbear-demo/pkg/workflow"
	"github.com/nite-coder/blackbear/pkg/log"
)

var (
	// grpc clients
	_walletClient walletProto.WalletServiceClient
	_eventClient  eventProto.EventServiceClient
)

func initialize(cfg config.Configuration) error {
	var err error

	cfg.InitLogger("worker")

	_eventClient, err = eventGRPC.NewClient(cfg)
	if err != nil {
		return err
	}

	_walletClient, err = walletGRPC.NewClient(cfg)
	if err != nil {
		return err
	}

	manager := workflow.Manager{
		Config:       cfg,
		WalletClient: _walletClient,
		EventClient:  _eventClient,
	}

	workflow.SetManager(&manager)

	log.Info("worker server is initialized")
	return nil
}
