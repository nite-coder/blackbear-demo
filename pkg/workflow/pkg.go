package workflow

import (
	"github.com/jasonsoft/starter/internal/pkg/config"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
)

var _manager *Manager

func SetManager(m *Manager) {
	_manager = m
}

type Manager struct {
	Config       config.Configuration
	WalletClient walletProto.WalletServiceClient
	EventClient  eventProto.EventServiceClient
}
