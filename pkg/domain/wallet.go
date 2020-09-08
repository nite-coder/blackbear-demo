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

// WalletServicer handles wallet's business logic
type WalletServicer interface {
	Wallet(ctx context.Context) (*Wallet, error)
	Withdraw(ctx context.Context, transID string, amount int64) error
}
