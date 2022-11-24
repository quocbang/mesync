package interceptors

import (
	"context"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryRecoveryInterceptor will recover system after panic behavior occurs.
func UnaryRecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "panic err: %v", e)
		}
	}()
	return handler(ctx, req)
}

// StreamRecoveryInterceptor will recover system after panic behavior occurs.
func StreamRecoveryInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "panic err: %v", e)
		}
	}()
	return handler(srv, stream)
}
