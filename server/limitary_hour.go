package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"gitlab.kenda.com.tw/kenda/mcom"
	mcomErr "gitlab.kenda.com.tw/kenda/mcom/errors"

	pb "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
)

// Create Limitary-hour
func (s Server) CreateLimitaryHour(ctx context.Context, req *pb.CreateLimitaryHourRequest) (*empty.Empty, error) {
	limitaryHours := parseLimitaryHour(req.GetLimitaryHour())
	if len(limitaryHours) == 0 {
		return nil, mcomErr.Error{
			Code:    mcomErr.Code_INSUFFICIENT_REQUEST,
			Details: "productType not found",
		}
	}

	return &empty.Empty{}, s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.CreateLimitaryHour(ctx, mcom.CreateLimitaryHourRequest{
			LimitaryHour: limitaryHours,
		})
	})
}

func parseLimitaryHour(limitaryHour []*pb.LimitaryHour) []mcom.LimitaryHour {
	limitaryHours := make([]mcom.LimitaryHour, len(limitaryHour))
	for i, v := range limitaryHour {
		if v == nil {
			continue
		}

		limitaryHours[i] = mcom.LimitaryHour{
			ProductType: v.GetProductType(),
			LimitaryHour: mcom.LimitaryHourParameter{
				Min: v.GetLimitaryHour().GetMin(),
				Max: v.GetLimitaryHour().GetMax(),
			},
		}
	}
	return limitaryHours
}
