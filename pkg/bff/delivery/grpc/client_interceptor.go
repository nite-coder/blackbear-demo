package gql

import (
	"context"
	"fmt"

	"github.com/jasonsoft/log/v2"
	internalMiddleware "github.com/jasonsoft/starter/internal/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ClientInterceptor return a client side interceptor for grpc
func ClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		logger := log.FromContext(ctx)

		// pass request_id to other services
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}
		md["request_id"] = []string{internalMiddleware.RequestIDFromContext(ctx)}

		//logger.Debugf("dump client metadata: %#v", md)

		// run
		err = invoker(metadata.NewOutgoingContext(ctx, md), method, req, resp, cc, opts...)

		// centralized error log
		if err != nil {
			logger.Interface("req", req).StackTrace().Err(err).Warnf("grpc client call failed, method: %s", method)
			err = fmt.Errorf("grpc client call failed, method: %s, %w", method, err)
		}

		return err
	}
}
