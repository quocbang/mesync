package interceptors

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	commonsCtx "gitlab.kenda.com.tw/kenda/commons/v2/utils/context"
)

func TestUnaryLoggingInterceptor(t *testing.T) {
	assert := assert.New(t)
	uli := UnaryLoggingInterceptor(getContextFields, contextWithLogger)
	{ // success
		_, err := uli(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/path/to/api",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.NoError(err)
	}
	{ // with error
		_, err := uli(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/path/to/api",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, errors.New("error after handler")
		})
		assert.Error(err)
	}

}

func getContextFields(_ context.Context) []zap.Field {
	return []zap.Field{
		zap.String("request_id", "XX001"),
		zap.String("user_id", "tester"),
		zap.Strings("factory_ids", []string{"Factory1", "Factory2"}),
	}
}

func contextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return commonsCtx.WithLogger(ctx, logger)
}
