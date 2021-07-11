package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nite-coder/blackbear-demo/internal/pkg/config"
	"github.com/nite-coder/blackbear-demo/pkg/domain"
	"github.com/nite-coder/blackbear-demo/pkg/wallet/proto"
	"github.com/nite-coder/blackbear/pkg/log"
)

// WalletServer is grpc server for wallet
type WalletServer struct {
	config        config.Configuration
	walletService domain.WalletUsecase
}

// NewWalletServer return a WalletServer instance
func NewWalletServer(cfg config.Configuration, walletService domain.WalletUsecase) *WalletServer {
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

	data, err := walletToGRPC(wallet)
	if err != nil {
		return nil, err
	}

	result := proto.GetWalletResponse{
		Data: data,
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
