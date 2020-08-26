package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/wallet"
	"github.com/jasonsoft/starter/pkg/wallet/proto"
)

type WalletServer struct {
	config        config.Configuration
	walletService wallet.WalletServicer
}

// NewWalletServer return a WalletServer instance
func NewWalletServer(cfg config.Configuration, walletService wallet.WalletServicer) *WalletServer {
	return &WalletServer{
		config:        cfg,
		walletService: walletService,
	}
}

// GetWallet returns single wallet instance
func (s *WalletServer) GetWallet(ctx context.Context, in *empty.Empty) (*proto.GetWalletResponse, error) {
	logger := log.FromContext(ctx)
	logger.Debug("grpc: begin GetWallet fn")

	//return nil, errors.New("oops..get wallet")

	wallet, err := s.walletService.Wallet(ctx)
	if err != nil {
		return nil, err
	}

	result := proto.GetWalletResponse{
		Data: walletToGWallet(wallet),
	}
	return &result, nil
}

// Withdraw fn will withdraw money from wallet
func (s *WalletServer) Withdraw(ctx context.Context, in *proto.WithdrawRequest) (*empty.Empty, error) {
	logger := log.FromContext(ctx)
	logger.Debug("grpc: begin Withdraw fn")

	err := s.walletService.Withdraw(ctx, in.TransId, in.Amount)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func walletToGWallet(wallet *wallet.Wallet) *proto.Wallet {
	updatedAt, _ := ptypes.TimestampProto(wallet.UpdatedAt)

	result := proto.Wallet{
		Id:        wallet.ID,
		Amount:    wallet.Amount,
		UpdatedAt: updatedAt,
	}

	return &result
}
