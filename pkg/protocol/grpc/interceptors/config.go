package interceptors

import (
	"context"
	"fmt"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	GenerateRequestID func() string
	ParseLogFields    func(context.Context) []zap.Field
	WithLoggerFunc    func(context.Context, *zap.Logger) context.Context

	Timeout time.Duration
	AuthKey string
}

// Build returns gRPC.Server config option that turn on unary/stream logging, recovery, some other interceptors.
func (conf Config) Build() ([]grpc.ServerOption, error) {
	if conf.GenerateRequestID == nil {
		return nil, fmt.Errorf("missing GenerateRequestID")
	}

	if conf.ParseLogFields == nil {
		return nil, fmt.Errorf("missing ParseLogFields")
	}

	if conf.WithLoggerFunc == nil {
		return nil, fmt.Errorf("missing WithLoggerFunc")
	}

	validateFunc := newValidateFunc(conf.AuthKey)

	return []grpc.ServerOption{
		// Add unary interceptor
		// interceptors.UnaryAuthenticationInterceptor and
		// interceptors.UnaryLoggingInterceptor depend on interceptors.UnaryContextInterceptor.
		// so, they must be call after interceptors.UnaryContextInterceptor.
		grpc_middleware.WithUnaryServerChain(
			UnaryRecoveryInterceptor,
			UnaryContextInterceptor(conf.GenerateRequestID),
			UnaryLoggingInterceptor(conf.ParseLogFields, conf.WithLoggerFunc),
			UnaryContextValidationInterceptor(validateFunc),
			UnaryTimeoutInterceptor(conf.Timeout),
		),
		// Add stream interceptor
		// interceptors.StreamLoggingInterceptor depend on interceptors.StreamContextInterceptor.
		// so, they must be call after interceptors.StreamContextInterceptor.
		grpc_middleware.WithStreamServerChain(
			StreamRecoveryInterceptor,
			StreamContextInterceptor(conf.GenerateRequestID),
			StreamLoggingInterceptor(conf.ParseLogFields, conf.WithLoggerFunc),
			StreamContextValidationInterceptor(validateFunc),
		),
	}, nil
}
