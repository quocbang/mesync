package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.kenda.com.tw/kenda/mcom"
	mcomErr "gitlab.kenda.com.tw/kenda/mcom/errors"
	"gitlab.kenda.com.tw/kenda/mcom/mock"

	"gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
)

func Test_CreateLimitaryHour(t *testing.T) {
	const (
		testProductType = "Type A"
		testMin         = 1
		testMax         = 10
	)

	assert := assert.New(t)
	ctx := contextWithFactoryIDs(context.Background(), factoryIDs...)

	{ // CreateLimitaryHour: insufficient request
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncCreateLimitaryHour,
			Input: mock.Input{
				Request: mcom.CreateLimitaryHourRequest{
					LimitaryHour: []mcom.LimitaryHour{},
				},
			},
			Output: mock.Output{
				Response: nil,
				Error: mcomErr.Error{
					Code: mcomErr.Code_INSUFFICIENT_REQUEST,
				},
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.CreateLimitaryHour(ctx, &mesync.CreateLimitaryHourRequest{})
		assert.Error(err)
		assert.Equal(mcomErr.Error{
			Code:    mcomErr.Code_INSUFFICIENT_REQUEST,
			Details: "productType not found",
		}, err)
	}
	{ // CreateLimitaryHour: good case
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncCreateLimitaryHour,
			Input: mock.Input{
				Request: mcom.CreateLimitaryHourRequest{
					LimitaryHour: []mcom.LimitaryHour{{
						ProductType: testProductType,
						LimitaryHour: mcom.LimitaryHourParameter{
							Min: testMin,
							Max: testMax,
						},
					}},
				},
			},
			Output: mock.Output{
				Response: nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.CreateLimitaryHour(ctx, &mesync.CreateLimitaryHourRequest{
			LimitaryHour: []*mesync.LimitaryHour{{
				ProductType: testProductType,
				LimitaryHour: &mesync.LimitaryHourParameter{
					Min: testMin,
					Max: testMax,
				},
			}},
		})
		assert.NoError(err)
	}
}
