package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"gitlab.kenda.com.tw/kenda/mcom"

	pb "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
)

// CreateStationGroup implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) CreateStationGroup(ctx context.Context, req *pb.CreateStationGroupRequest) (*empty.Empty, error) {
	err := s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.CreateStationGroup(ctx, mcom.StationGroupRequest{
			ID:       req.GetId(),
			Stations: req.GetInfo().GetStations(),
		})
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// UpdateStationGroup implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) UpdateStationGroup(ctx context.Context, req *pb.UpdateStationGroupRequest) (*empty.Empty, error) {
	err := s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.UpdateStationGroup(ctx, mcom.StationGroupRequest{
			ID:       req.GetId(),
			Stations: req.GetInfo().GetStations(),
		})
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// DeleteStationGroup implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) DeleteStationGroup(ctx context.Context, req *pb.DeleteStationGroupRequest) (*empty.Empty, error) {
	err := s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.DeleteStationGroup(ctx, mcom.DeleteStationGroupRequest{
			GroupID: req.GetId(),
		})
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
