package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// UnaryTimeoutInterceptor checks if handler request is timeout.
// @Param: timeout which set the context's timeout.
func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel() // releases resources if handler completes before timeout elapses

		return handler(newCtx, req)
	}
}

// StreamTimeoutInterceptor checks if handler request is timeout.
// @Param: timeout which set the context's timeout.
func StreamTimeoutInterceptor(timeout time.Duration) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		stream := serverStream{
			ServerStream: ss,
			ctx:          ss.Context(),
		}

		newCtx, cancel := context.WithTimeout(stream.Context(), timeout)
		stream.ctx = newCtx
		defer cancel() // releases resources if handler completes before timeout elapses

		return handler(srv, stream)
	}
}
