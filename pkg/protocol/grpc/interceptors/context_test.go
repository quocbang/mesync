package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	grpc_context "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/context"
	mesyncMetadata "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/metadata"
)

const (
	authKey   = "a12345"
	requestID = "XXXXXX123456"
	userID    = "tester"

	defaultGeneratedRID = "ABCDEF12345"
)

func TestUnaryContextInterceptor(t *testing.T) {
	assert := assert.New(t)
	{ // without metadata assigned
		ctx := context.Background()
		uci := UnaryContextInterceptor(defaultRequestID)
		// nolint
		uci(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(defaultGeneratedRID, grpc_context.GetRequestID(ctx))
			assert.Empty(grpc_context.GetAuthKey(ctx))
			assert.Empty(grpc_context.GetUserID(ctx))
			return nil, nil
		})
	}
	{ // with both requestID and authKey assigned on metadata
		ctx := metadata.NewIncomingContext(context.Background(),
			metadata.Pairs(
				mesyncMetadata.RequestID, requestID,
				mesyncMetadata.AuthenticationKey, authKey,
				mesyncMetadata.UserID, userID,
			),
		)
		uci := UnaryContextInterceptor(defaultRequestID)
		// nolint
		uci(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(requestID, grpc_context.GetRequestID(ctx))
			assert.Equal(authKey, grpc_context.GetAuthKey(ctx))
			assert.Equal(userID, grpc_context.GetUserID(ctx))
			return nil, nil
		})
	}
	{ // only authKey assigned and generate a default requestID
		ctx := metadata.NewIncomingContext(context.Background(),
			metadata.Pairs(
				mesyncMetadata.AuthenticationKey, authKey,
				mesyncMetadata.UserID, userID,
			),
		)
		uci := UnaryContextInterceptor(defaultRequestID)
		// nolint
		uci(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(defaultGeneratedRID, grpc_context.GetRequestID(ctx))
			assert.Equal(authKey, grpc_context.GetAuthKey(ctx))
			assert.Equal(userID, grpc_context.GetUserID(ctx))
			return nil, nil
		})
	}
}

func defaultRequestID() string {
	return defaultGeneratedRID
}
