package server

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	pbTypes "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/types"
)

func Test_toPointerDecimal(t *testing.T) {
	assert := assert.New(t)

	// with value
	{
		const v = "123.5"

		d := &pbTypes.Decimal{
			Value: v,
		}
		n, err := toPointerDecimal(d)
		if assert.NoError(err) {
			expected := decimal.RequireFromString(v)
			assert.Equal(&expected, n)
		}
	}
	// nil value
	{
		var d *pbTypes.Decimal
		assert.Nil(toPointerDecimal(d))
	}
}
