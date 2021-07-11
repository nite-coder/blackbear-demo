package workflow

import (
	"go.temporal.io/sdk/temporal"
	"google.golang.org/grpc/status"
)

func centralizedError(err error) error {

	// grpc
	grpcStatus, ok := status.FromError(err)
	if ok {
		errType := "GRPC_" + grpcStatus.Code().String()
		return temporal.NewApplicationError(grpcStatus.Message(), errType)
	}

	return err
}
