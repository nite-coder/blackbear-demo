package service

import (
	"context"
	"sync"
	"time"

	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/wallet"
)

type WalletService struct {
	mu     sync.Mutex
	config config.Configuration
	wallet *wallet.Wallet
}

func NewWalletService(cfg config.Configuration) wallet.WalletServicer {
	return &WalletService{
		config: cfg,
		wallet: &wallet.Wallet{
			ID:        1,
			Amount:    10000,
			UpdatedAt: time.Now().UTC(),
		},
	}
}

func (svc *WalletService) Wallet(ctx context.Context) (*wallet.Wallet, error) {
	return svc.wallet, nil
}

func (svc *WalletService) Withdraw(ctx context.Context, transID string, amount int64) error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.wallet.Amount -= amount
	svc.wallet.UpdatedAt = time.Now().UTC()

	return nil
}
