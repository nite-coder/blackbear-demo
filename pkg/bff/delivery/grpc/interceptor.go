package gql

import (
	"context"

	internalMiddleware "github.com/jasonsoft/starter/internal/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ClientInterceptor return a client side interceptor for grpc
func ClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		// logger := log.FromContext(ctx)
		// logger.Debug("== begin client invoker ==")

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}

		md["request_id"] = []string{internalMiddleware.RequestIDFromContext(ctx)}

		err = invoker(metadata.NewOutgoingContext(ctx, md), method, req, resp, cc, opts...)
		//logger.Debug("== end client invoker ==")

		return nil
	}
}
