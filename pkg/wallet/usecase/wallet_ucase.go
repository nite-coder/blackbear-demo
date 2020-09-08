package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/domain"
)

type walletUsecase struct {
	mu     sync.Mutex
	config config.Configuration
	wallet *domain.Wallet
}

// NewWalletUsecase create a new walletUsecase object representation of domain.WalletUsecase interface
func NewWalletUsecase(cfg config.Configuration) domain.WalletUsecase {
	return &walletUsecase{
		config: cfg,
		wallet: &domain.Wallet{
			ID:        1,
			Amount:    10000,
			UpdatedAt: time.Now().UTC(),
		},
	}
}

// Wallet returns a wallet
func (u *walletUsecase) Wallet(ctx context.Context) (*domain.Wallet, error) {
	logger := log.FromContext(ctx)
	logger.Debug("begin wallet fn")

	return u.wallet, nil
}

// Withdraw fn remove amount from wallet
func (u *walletUsecase) Withdraw(ctx context.Context, transID string, amount int64) error {
	logger := log.FromContext(ctx)
	logger.Debug("begin withdraw fn")

	u.mu.Lock()
	defer u.mu.Unlock()

	u.wallet.Amount -= amount
	u.wallet.UpdatedAt = time.Now().UTC()

	return nil
}
