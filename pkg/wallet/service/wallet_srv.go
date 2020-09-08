package service

import (
	"context"
	"sync"
	"time"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/domain"
)

// WalletService handles wallet's business logic
type WalletService struct {
	mu     sync.Mutex
	config config.Configuration
	wallet *domain.Wallet
}

// NewWalletService create an instance of wallet service
func NewWalletService(cfg config.Configuration) domain.WalletServicer {
	return &WalletService{
		config: cfg,
		wallet: &domain.Wallet{
			ID:        1,
			Amount:    10000,
			UpdatedAt: time.Now().UTC(),
		},
	}
}

// Wallet returns a wallet
func (svc *WalletService) Wallet(ctx context.Context) (*domain.Wallet, error) {
	logger := log.FromContext(ctx)
	logger.Debug("begin wallet fn")

	return svc.wallet, nil
}

// Withdraw fn remove amount from wallet
func (svc *WalletService) Withdraw(ctx context.Context, transID string, amount int64) error {
	logger := log.FromContext(ctx)
	logger.Debug("begin withdraw fn")

	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.wallet.Amount -= amount
	svc.wallet.UpdatedAt = time.Now().UTC()

	return nil
}
