package server

import (
	"context"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"

	"gitlab.kenda.com.tw/kenda/mcom"
	"gitlab.kenda.com.tw/kenda/mcom/utils/types"

	"gitlab.kenda.com.tw/kenda/mesync/pkg/config"
	"gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
	grpcContext "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/context"
)

func TestServer_Cloud(t *testing.T) {
	assert := assert.New(t)
	server := CloudServer{
		cloud: cloud{
			configs: config.CloudConfigs{
				FactoryContainerMap: map[string]string{"mestest": "mestest"},
				ContainerToken: map[string]config.BlobToken{
					"mestest": {
						AccountName: "accountName",
						SasToken:    "token",
					},
				},
			},
		},
	}
	timeNow := time.Now()
	timestampNow, err := ptypes.TimestampProto(timeNow)
	assert.NoError(err)

	pHttpReq := monkey.Patch(sendHttpRequest, func(rub reUploadBlob) error {
		return nil
	})
	req := mesync.UploadBlobRequest{
		ParametersHeader: []string{"greeting"},
		Detail: []*mesync.PerUploadBlobRequest{
			{
				DateTime: timestampNow,
				ManufactureParameters: &mesync.ManufactureParameters{
					Value: []*mesync.ManufactureParametersValue{
						{
							OneofName: &mesync.ManufactureParametersValue_StringValue{StringValue: "hello"},
						},
					},
				},
			},
		},
		StationId: "123",
		Batch:     456,
		Resources: []string{"b1"},
	}
	ctx := context.WithValue(context.Background(), grpcContext.FactoryIDs, []string{"mestest"})
	mcomReq, err := server.uploadBlob(ctx, &req)
	assert.NoError(err)

	uriTime, err := ptypes.Timestamp(timestampNow)
	assert.NoError(err)

	assert.Equal(mcom.CreateBlobResourceRecordRequest{
		Details: []mcom.CreateBlobResourceRecordDetail{{
			BlobURI:       uriTime.Format("2006/01/02/") + sha256String("123") + uriTime.Format("/150405"),
			Resources:     []string{"b1"},
			Station:       "123",
			DateTime:      types.ToTimeNano(timeNow),
			ContainerName: "mestest",
		}},
	}, mcomReq)
	pHttpReq.Unpatch()
}
