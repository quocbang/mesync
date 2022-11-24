package server

import (
	"github.com/shopspring/decimal"

	"gitlab.kenda.com.tw/kenda/mcom/utils/types"

	pbTypes "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/types"
)

// toPointerDecimal returns the pointer of decimal.Decimal type, if the
// d value is nil, toPointerDecimal returns nil too.
func toPointerDecimal(d *pbTypes.Decimal) (*decimal.Decimal, error) {
	if d == nil {
		return nil, nil
	}

	return types.Decimal.NewFromString(d.GetValue())
}
