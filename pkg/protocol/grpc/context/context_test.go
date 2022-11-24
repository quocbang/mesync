package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	mesyncMetadata "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/metadata"
)

const (
	reqID          = "XXX001"
	authKey        = "a12345"
	generatedReqID = "AAXXBB"
	userID         = "TESTER"
)

var factoryIDs = []string{"KY", "KU"}

func TestParser_Parse(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	{ // without field options
		md := metadata.Pairs(
			mesyncMetadata.RequestID, reqID,
		)
		p := NewParser(ctx, md)
		pCtx := p.Parse(mesyncMetadata.RequestID, RequestID).Done()
		assert.Equal(context.WithValue(ctx, RequestID, reqID), pCtx)
	}
	{ // with field options
		p := NewParser(ctx, nil)
		pCtx := p.Parse(mesyncMetadata.RequestID, RequestID, WithDefault(func() string {
			return generatedReqID
		})).Done()
		assert.Equal(context.WithValue(ctx, RequestID, generatedReqID), pCtx)
	}
	{ // get md fields
		md := metadata.Pairs(
			mesyncMetadata.RequestID, reqID,
			mesyncMetadata.AuthenticationKey, authKey,
			mesyncMetadata.UserID, userID,
		)
		md.Append(mesyncMetadata.FactoryIDs, factoryIDs...)

		p := NewParser(ctx, md)
		pCtx := p.Parse(mesyncMetadata.RequestID, RequestID).
			Parse(mesyncMetadata.AuthenticationKey, AuthKey).
			Parse(mesyncMetadata.UserID, UserID).
			Parse(mesyncMetadata.FactoryIDs, FactoryIDs, WithMultiple()).
			Done()
		assert.Equal(reqID, GetRequestID(pCtx))
		assert.Equal(authKey, GetAuthKey(pCtx))
		assert.Equal(factoryIDs, GetFactoryIDs(pCtx))
	}
	{ // get md fields and factoryIDs without Multiple option
		md := metadata.Pairs(
			mesyncMetadata.RequestID, reqID,
			mesyncMetadata.AuthenticationKey, authKey,
			mesyncMetadata.UserID, userID,
		)
		md.Append(mesyncMetadata.FactoryIDs, factoryIDs...)

		p := NewParser(ctx, md)
		pCtx := p.Parse(mesyncMetadata.RequestID, RequestID).
			Parse(mesyncMetadata.AuthenticationKey, AuthKey).
			Parse(mesyncMetadata.UserID, UserID).
			Parse(mesyncMetadata.FactoryIDs, FactoryIDs).
			Done()
		assert.Equal(reqID, GetRequestID(pCtx))
		assert.Equal(authKey, GetAuthKey(pCtx))
		assert.Equal([]string{factoryIDs[0]}, GetFactoryIDs(pCtx))
	}
}

func Test_getStringsContext(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	md := metadata.Pairs()
	md.Append(mesyncMetadata.FactoryIDs, factoryIDs...)
	p := NewParser(ctx, md)

	{ // return empty struct if context not contant factoryID key
		ftyIDs := getStringsContext(ctx, FactoryIDs)
		assert.Equal([]string{}, ftyIDs)
	}
	{ // return multiple values if context multiple values for the same key
		pCtx := p.Parse(mesyncMetadata.FactoryIDs, FactoryIDs, WithMultiple()).Done()
		ftyIDs := getStringsContext(pCtx, FactoryIDs)
		assert.Equal(factoryIDs, ftyIDs)
	}
	{ // single value
		pCtx := p.Parse(mesyncMetadata.FactoryIDs, FactoryIDs).Done()
		ftyIDs := getStringsContext(pCtx, FactoryIDs)
		assert.Equal([]string{factoryIDs[0]}, ftyIDs)
	}
}
