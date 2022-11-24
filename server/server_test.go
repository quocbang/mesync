package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	commonsCtx "gitlab.kenda.com.tw/kenda/commons/v2/utils/context"
	"gitlab.kenda.com.tw/kenda/mcom"
)

type emptyDataManager struct {
	mcom.DataManager
}

func TestServer_getDataManagerByFactory(t *testing.T) {
	assert := assert.New(t)

	const (
		keyA = "A"
		keyB = "B"
		keyC = "C"

		keyNotFound = "୧༼◕ ᴥ ◕༽୨"
	)

	s := NewServer(ServerParams{
		LocalDataManager: map[string]mcom.DataManager{
			keyA: new(emptyDataManager),
			keyB: new(emptyDataManager),
			keyC: new(emptyDataManager),
		},
		Logger: func(ctx context.Context) *zap.Logger { return commonsCtx.Logger(ctx) }})

	// found the factory
	dm, err := s.getDataManagerByFactory(keyB)
	assert.NoError(err)
	assert.NotNil(dm)

	// the factory not found
	dm, err = s.getDataManagerByFactory(keyNotFound)
	assert.EqualError(err, fmt.Sprintf("mcom service not found [%s]", keyNotFound))
	assert.Nil(dm)
}

func TestServer_eachFactory(t *testing.T) {
	assert := assert.New(t)

	const (
		keyA = "A"
		keyB = "B"
		keyC = "C"
	)

	s := NewServer(ServerParams{
		LocalDataManager: map[string]mcom.DataManager{
			keyA: new(emptyDataManager),
			keyB: new(emptyDataManager),
			keyC: new(emptyDataManager),
		},
		Logger: func(ctx context.Context) *zap.Logger { return commonsCtx.Logger(ctx) }})

	// no error occurred
	{
		handler := func(ctx context.Context, factoryID string, i int, dm mcom.DataManager) error {
			return nil
		}

		ctx := contextWithFactoryIDs(context.Background(), keyA, keyB, keyC)
		assert.NoError(s.eachFactory(ctx, handler))
	}
	// occurred error
	{
		handler := func(ctx context.Context, factoryID string, i int, dm mcom.DataManager) error {
			return fmt.Errorf("↑_(ΦwΦ;)Ψ")
		}

		ctx := contextWithFactoryIDs(context.Background(), keyA, keyB, keyC)
		assert.Error(s.eachFactory(ctx, handler))
	}
	// there is no factory id in the context
	{
		handler := func(ctx context.Context, factoryID string, i int, dm mcom.DataManager) error {
			return nil
		}

		assert.EqualError(s.eachFactory(context.Background(), handler), "missing factory id")
	}
}

func Test_mergeErrors(t *testing.T) {
	assert := assert.New(t)

	// there are some errors
	{
		const (
			errA = "(´･ω･`)"
			errB = "(ﾉД`)ｼｸｼｸ"
		)

		ch := make(chan error, 3)
		ch <- fmt.Errorf(errA)
		ch <- fmt.Errorf(errB)
		close(ch)

		assert.EqualError(mergeErrors(ch), fmt.Sprintf("%s; %s", errA, errB))
	}
	// there is no error
	{
		ch := make(chan error, 3)
		close(ch)

		assert.NoError(mergeErrors(ch))
	}
}
