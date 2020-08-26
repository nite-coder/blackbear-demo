package grpc

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/jasonsoft/starter/pkg/wallet"
	"github.com/jasonsoft/starter/pkg/wallet/proto"
)

func walletToGRPC(wallet *wallet.Wallet) (*proto.Wallet, error) {
	if wallet == nil {
		return nil, nil
	}

	updatedAt, err := ptypes.TimestampProto(wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}

	result := proto.Wallet{
		Id:        wallet.ID,
		Amount:    wallet.Amount,
		UpdatedAt: updatedAt,
	}

	return &result, nil
}
