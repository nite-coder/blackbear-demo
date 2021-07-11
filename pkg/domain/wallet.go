package domain

import (
	"context"
	"time"
)

// Wallet is Wallet
type Wallet struct {
	ID        int64
	Amount    int64
	UpdatedAt time.Time
}

// WalletUsecase represents the wallet's usecases
type WalletUsecase interface {
	Wallet(ctx context.Context) (*Wallet, error)
	Withdraw(ctx context.Context, transID string, amount int64) error
}
