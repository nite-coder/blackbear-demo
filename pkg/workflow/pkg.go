package workflow

import (
	"context"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	internalMiddleware "github.com/jasonsoft/starter/internal/pkg/middleware"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
)

var _manager *Manager

// SetManager allow extenal to pass Manager instance to workflow.  The manager instance will be used in the package.
func SetManager(m *Manager) {
	_manager = m
}

// Manager is a group of grpc clients.
type Manager struct {
	Config       config.Configuration
	WalletClient walletProto.WalletServiceClient
	EventClient  eventProto.EventServiceClient
}

func getLogger(ctx context.Context) log.Context {
	requestID := internalMiddleware.RequestIDFromContext(ctx)

	if len(requestID) > 0 {
		return log.Str("request_id", requestID)
	}

	return log.FromContext(ctx)
}
