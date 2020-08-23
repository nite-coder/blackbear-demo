package service

import (
	"context"
	"sync"

	"github.com/jasonsoft/starter/pkg/wallet"
)

type WalletService struct {
	mu     sync.Mutex
	wallet *wallet.Wallet
}

func NewWalletService(wallet *wallet.Wallet) wallet.WalletServicer {
	return &WalletService{
		wallet: wallet,
	}
}

func (svc *WalletService) Wallet(ctx context.Context) (*wallet.Wallet, error) {
	return svc.wallet, nil
}

func (svc *WalletService) Withdraw(ctx context.Context, transID string, amount int64) error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.wallet.Amount -= amount

	return nil
}
