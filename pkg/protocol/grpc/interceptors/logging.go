package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ContextWithLogger func(context.Context, *zap.Logger) context.Context

// UnaryLoggingInterceptor logs gRPC request and response message with tracking ID.
// Also, it is associated with a logger with log fields.
// @Param: getContextFields to log corresponding header field after parsing metadata context using UnaryContextInterceptor.
func UnaryLoggingInterceptor(getContextFields func(context.Context) []zap.Field, contextWithLogger ContextWithLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger := zap.L().With(
			append([]zap.Field{zap.String("path", info.FullMethod)}, getContextFields(ctx)...)...,
		)
		newCtx := contextWithLogger(ctx, logger)

		logger.Debug("start unary call", zap.Any("request", req))

		now := time.Now()
		resp, err := handler(newCtx, req)
		if err != nil {
			logger.Error("finished unary call with error",
				zap.Any("request", req),
				zap.Any("response", resp),
				zap.Duration("elapsed_time", time.Since(now)),
				zap.Error(err),
			)
			return resp, err
		}
		logger.Debug("finished unary call",
			zap.Any("response", resp),
			zap.Duration("elapsed_time", time.Since(now)))
		return resp, err
	}
}

// StreamLoggingInterceptor logs gRPC request and response message with tracking ID.
// Also, it is associated with a logger with log fields.
// @Param: getContextFields to log corresponding header field after parsing metadata context using UnaryContextInterceptor.
func StreamLoggingInterceptor(getContextFields func(context.Context) []zap.Field, contextWithLogger ContextWithLogger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		stream := serverStream{
			ServerStream: ss,
			ctx:          ss.Context(),
		}

		ctx := stream.Context()

		logger := zap.L().With(
			append([]zap.Field{zap.String("path", info.FullMethod)}, getContextFields(ctx)...)...,
		)
		stream.ctx = contextWithLogger(ctx, logger)

		now := time.Now()
		err := handler(srv, stream)
		if err != nil {
			logger.Error("finished stream call with error",
				zap.Any("request", srv),
				zap.Duration("elapsed_time", time.Since(now)),
				zap.Error(err),
			)
			return err
		}
		logger.Debug("finished stream call",
			zap.Duration("elapsed_time", time.Since(now)))
		return err
	}
}
