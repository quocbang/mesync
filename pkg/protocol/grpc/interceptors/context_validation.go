package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_context "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/context"
)

// UnaryContextValidationInterceptor validates the incoming context condition.
func UnaryContextValidationInterceptor(validate func(context.Context) error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := validate(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// StreamContextValidationInterceptor validates the incoming context condition.
func StreamContextValidationInterceptor(validate func(context.Context) error) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := validate(ss.Context()); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

// newValidateFunc creates a new validation function to validate
// the user id and the auth key from the context.
func newValidateFunc(key string) func(context.Context) error {
	return func(ctx context.Context) error {
		// User-ID context must not empty
		if grpc_context.GetUserID(ctx) == "" {
			return status.Error(codes.FailedPrecondition, "empty User-ID context")
		}
		// request's authentication key must be equal to configuration's.
		if grpc_context.GetAuthKey(ctx) != key {
			return status.Error(codes.Unauthenticated, "mismatch authentication key")
		}
		return nil
	}
}
