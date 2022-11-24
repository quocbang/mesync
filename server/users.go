package server

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"

	"gitlab.kenda.com.tw/kenda/mcom"

	pb "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
)

func getUsersInfo(usersInfo []*pb.UserInfo) []mcom.User {
	users := make([]mcom.User, len(usersInfo))
	for i, v := range usersInfo {
		users[i] = mcom.User{
			ID:           v.GetId(),
			Account:      v.GetAccount(),
			DepartmentID: v.GetDepartmentId(),
		}
	}
	return users
}

// CreateUsers implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) CreateUsers(ctx context.Context, req *pb.CreateUsersRequest) (*empty.Empty, error) {
	err := s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.CreateUsers(ctx, mcom.CreateUsersRequest{
			Users: getUsersInfo(req.Users),
		})
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// UpdateUser implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*empty.Empty, error) {
	var leaveDate time.Time
	if req.GetLeaveDate() != nil {
		date, err := ptypes.Timestamp(req.GetLeaveDate())
		if err != nil {
			return nil, err
		}
		leaveDate = date.Local()
	}

	err := s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.UpdateUser(ctx, mcom.UpdateUserRequest{
			ID:           req.GetId(),
			Account:      req.GetAccount(),
			DepartmentID: req.GetDepartmentId(),
			LeaveDate:    leaveDate,
		})
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// DeleteUser implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*empty.Empty, error) {
	err := s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.DeleteUser(ctx, mcom.DeleteUserRequest{
			ID: req.GetId(),
		})
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// CreateDepartments implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) CreateDepartments(ctx context.Context, req *pb.Departments) (*empty.Empty, error) {
	err := s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.CreateDepartments(ctx, req.Ids)
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// UpdateDepartment implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) UpdateDepartment(ctx context.Context, req *pb.UpdateDepartmentRequest) (*empty.Empty, error) {
	err := s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.UpdateDepartment(ctx, mcom.UpdateDepartmentRequest{
			OldID: req.GetId(),
			NewID: req.GetNewId(),
		})
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// DeleteDepartment implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) DeleteDepartment(ctx context.Context, req *pb.DeleteDepartmentRequest) (*empty.Empty, error) {
	err := s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.DeleteDepartment(ctx, mcom.DeleteDepartmentRequest{
			DepartmentID: req.GetId(),
		})
	})
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
