package server

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/shopspring/decimal"

	"gitlab.kenda.com.tw/kenda/mcom"

	pb "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
)

// CreateProductPlan implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) CreateProductPlan(ctx context.Context, req *pb.CreatePlanRequest) (*empty.Empty, error) {
	t, err := ptypes.Timestamp(req.GetDate())
	if err != nil {
		return nil, err
	}

	quantity, err := decimal.NewFromString(req.GetQuantity().GetValue())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		existed, err := dm.IsProductExisted(ctx, req.GetProduct().GetId())
		if err != nil {
			return err
		}
		// internal library error
		if !existed {
			return fmt.Errorf("product=%s has not been released / existed", req.GetProduct().GetId())
		}
		return dm.CreateProductPlan(ctx, mcom.CreateProductionPlanRequest{
			Date: t,
			Product: mcom.Product{
				ID:   req.GetProduct().GetId(),
				Type: req.GetProduct().GetType(),
			},
			DepartmentOID: req.GetDepartmentId(),
			Quantity:      quantity,
		})
	})
}
