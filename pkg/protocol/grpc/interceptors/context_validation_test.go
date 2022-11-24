package interceptors

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_context "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/context"
)

func TestUnaryUserInterceptor(t *testing.T) {
	assert := assert.New(t)
	{ // normal
		ucvi := UnaryContextValidationInterceptor(func(_ context.Context) error {
			return nil
		})
		_, err := ucvi(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.NoError(err)
	}
	{ // error
		ucvi := UnaryContextValidationInterceptor(func(_ context.Context) error {
			return errors.New("unauthenticated")
		})
		_, err := ucvi(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.EqualError(err, "unauthenticated")
	}
}

func Test_newValidateFunc(t *testing.T) {
	assert := assert.New(t)

	const (
		userID  = "tester"
		authKey = "a12345"
	)

	var (
		withUserID = func(ctx context.Context) context.Context {
			return context.WithValue(ctx, grpc_context.UserID, userID)
		}
		withAuthKey = func(ctx context.Context) context.Context {
			return context.WithValue(ctx, grpc_context.AuthKey, authKey)
		}
	)

	{ // normal
		ctx := withUserID(context.Background())
		ctx = withAuthKey(ctx)

		assert.Nil(newValidateFunc(authKey)(ctx))
	}
	{ // empty user id

		assert.Equal(newValidateFunc(authKey)(withAuthKey(context.Background())), status.Error(codes.FailedPrecondition, "empty User-ID context"))
	}
	{ // mismatch authentication key
		ctx := withUserID(context.Background())
		ctx = withAuthKey(ctx)

		assert.Equal(newValidateFunc("ab")(ctx), status.Error(codes.Unauthenticated, "mismatch authentication key"))
	}
}
