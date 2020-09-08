package worker

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	eventGRPC "github.com/jasonsoft/starter/pkg/event/delivery/grpc"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	walletGRPC "github.com/jasonsoft/starter/pkg/wallet/delivery/grpc"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
	"github.com/jasonsoft/starter/pkg/workflow"
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
