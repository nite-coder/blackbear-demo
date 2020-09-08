package grpc

import (
	"context"
	"errors"

	"github.com/jasonsoft/log/v2"
	internalMiddleware "github.com/jasonsoft/starter/internal/pkg/middleware"
	"github.com/jasonsoft/starter/pkg/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Interceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		logger := log.FromContext(ctx)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, grpc.Errorf(codes.DataLoss, "metadata is not found")
		}

		// get requestID from metadata and create a new log context
		var requestID string
		if val, ok := md["request_id"]; ok {
			requestID = val[0]
		}
		logger = logger.Str("request_id", requestID)
		ctx = logger.WithContext(ctx)
		ctx = internalMiddleware.SetRequestIDToContext(ctx, requestID)

		result, err := handler(ctx, req)

		// centralized error
		if err != nil {
			var appErr *domain.AppError
			if errors.As(err, &appErr) {
				gErr := status.Error(appErr.Status, appErr.Message)
				return result, gErr
			}

			// unknow error
			logger.Err(err).Error("wallet grpc unknown error")
			gErr := status.Error(codes.Unknown, err.Error())
			return result, gErr
		}

		return result, err
	}
}
