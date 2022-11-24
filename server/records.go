package server

import (
	"context"

	"github.com/golang/protobuf/ptypes"

	"gitlab.kenda.com.tw/kenda/mcom"

	pb "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
	pbTypes "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/types"
	grpc_context "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/context"
)

func composeFeedInfo(feedRecords []mcom.FeedRecord) []*pb.FeedRecord {
	listInfo := make([]*pb.FeedRecord, len(feedRecords))
	for i, res := range feedRecords {
		feedMaterials := make([]*pb.FeedMaterials, len(res.Materials))
		for mtrlIdx, mtrl := range res.Materials {
			feedMaterials[mtrlIdx] = &pb.FeedMaterials{
				Id:         mtrl.ID,
				Grade:      mtrl.Grade,
				ResourceId: mtrl.ResourceID,
				Quantity: &pbTypes.Decimal{
					Value: mtrl.Quantity.String(),
				},
				Site: &pb.Site{
					Name:  mtrl.Site.Name,
					Index: int32(mtrl.Site.Index),
				},
			}
		}

		listInfo[i] = &pb.FeedRecord{
			WorkOrder: res.WorkOrder,
			RecipeId:  res.RecipeID,
			Batch:     res.Batch,
			StationId: res.StationID,
			Materials: feedMaterials,
		}
	}
	return listInfo
}

// ListFeedRecords implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) ListFeedRecords(ctx context.Context, req *pb.ListRecordsRequest) (*pb.ListFeedRecordsReply, error) {
	ids := grpc_context.GetFactoryIDs(ctx)
	listRecords := make([]*pb.FactoryFeedRecord, len(ids))

	date, err := ptypes.Timestamp(req.GetDate())
	if err != nil {
		return nil, err
	}

	err = s.eachFactory(ctx, func(ctx context.Context, id string, index int, dm mcom.DataManager) error {
		results, err := dm.ListFeedRecords(ctx, mcom.ListRecordsRequest{
			Date:         date,
			DepartmentID: req.GetDepartmentId(),
		})
		if err != nil {
			return err
		}

		listRecords[index] = &pb.FactoryFeedRecord{
			FactoryId: id,
			Records:   composeFeedInfo(results),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.ListFeedRecordsReply{
		Records: listRecords,
	}, nil
}

func composeCollectInfo(collectRecords []mcom.CollectRecord) []*pb.CollectRecord {
	listInfo := make([]*pb.CollectRecord, len(collectRecords))
	for i, res := range collectRecords {
		listInfo[i] = &pb.CollectRecord{
			ResourceId: res.ResourceID,
			WorkOrder:  res.WorkOrder,
			RecipeId:   res.RecipeID,
			ProductId:  res.ProductID,
			Quantity: &pbTypes.Decimal{
				Value: res.Quantity.String(),
			},
		}
	}
	return listInfo
}

// ListCollectRecords implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) ListCollectRecords(ctx context.Context, req *pb.ListRecordsRequest) (*pb.ListCollectRecordsReply, error) {
	ids := grpc_context.GetFactoryIDs(ctx)
	listRecords := make([]*pb.FactoryCollectRecord, len(ids))

	date, err := ptypes.Timestamp(req.GetDate())
	if err != nil {
		return nil, err
	}

	err = s.eachFactory(ctx, func(ctx context.Context, id string, index int, dm mcom.DataManager) error {
		results, err := dm.ListCollectRecords(ctx, mcom.ListRecordsRequest{
			Date:         date,
			DepartmentID: req.GetDepartmentId(),
		})
		if err != nil {
			return err
		}

		listRecords[index] = &pb.FactoryCollectRecord{
			FactoryId: id,
			Records:   composeCollectInfo(results),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.ListCollectRecordsReply{
		Records: listRecords,
	}, nil
}
