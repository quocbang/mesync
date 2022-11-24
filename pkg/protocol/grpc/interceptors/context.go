package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	grpc_context "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/context"
	mesyncMetadata "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/metadata"
)

// UnaryContextInterceptor populates context parsed from MD.
func UnaryContextInterceptor(defaultRequestID func() string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		ctx = grpc_context.NewParser(ctx, md).
			Parse(mesyncMetadata.RequestID, grpc_context.RequestID, grpc_context.WithDefault(defaultRequestID)).
			Parse(mesyncMetadata.AuthenticationKey, grpc_context.AuthKey).
			Parse(mesyncMetadata.UserID, grpc_context.UserID).
			Parse(mesyncMetadata.FactoryIDs, grpc_context.FactoryIDs, grpc_context.WithMultiple()).
			Done()

		return handler(ctx, req)
	}
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss serverStream) Context() context.Context {
	return ss.ctx
}

func StreamContextInterceptor(defaultRequestID func() string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		stream := serverStream{
			ServerStream: ss,
			ctx:          ss.Context(),
		}

		md, _ := metadata.FromIncomingContext(stream.ctx)

		stream.ctx = grpc_context.NewParser(stream.ctx, md).
			Parse(mesyncMetadata.RequestID, grpc_context.RequestID, grpc_context.WithDefault(defaultRequestID)).
			Parse(mesyncMetadata.AuthenticationKey, grpc_context.AuthKey).
			Parse(mesyncMetadata.UserID, grpc_context.UserID).
			Parse(mesyncMetadata.FactoryIDs, grpc_context.FactoryIDs, grpc_context.WithMultiple()).
			Done()

		return handler(srv, stream)
	}
}
