package interceptors

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnaryTimeoutInterceptor(t *testing.T) {
	assert := assert.New(t)

	const (
		timeout = time.Millisecond * 200
		okReply = "ok"
	)
	var (
		interceptor = UnaryTimeoutInterceptor(timeout)
	)

	// no timeout
	{
		var handler = func(ctx context.Context, req interface{}) (interface{}, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return okReply, nil
			}
		}

		reply, err := interceptor(context.Background(), struct{}{}, nil, handler)
		assert.Equal(okReply, reply)
		assert.NoError(err)
	}
}
