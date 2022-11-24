package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnaryRecoveryInterceptor(t *testing.T) {
	assert := assert.New(t)

	assert.NotPanics(func() {
		reply, err := UnaryRecoveryInterceptor(context.Background(), struct{}{}, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			panic("failed")
		})
		assert.Empty(reply)
		assert.Error(err)
	})
}
