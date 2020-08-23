package wallet

import (
	"context"
	"time"
)

type Wallet struct {
	ID        int64
	Amount    int64
	UpdatedAt time.Time
}

type WalletServicer interface {
	Wallet(ctx context.Context) (*Wallet, error)
	Withdraw(ctx context.Context, transID string, amount int64) error
}
