package server

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"

	"gitlab.kenda.com.tw/kenda/mcom"
	mcomErr "gitlab.kenda.com.tw/kenda/mcom/errors"
	"gitlab.kenda.com.tw/kenda/mcom/mock"

	"gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
)

const (
	testUserA = "tester_A"
	testUserB = "tester_B"

	testADUser       = "AD_user"
	testADAccount    = "AD_account"
	testADNewAccount = "new_account"

	testDepartmentA = "dep_A"
	testDepartmentB = "dep_B"
)

var factoryIDs = []string{testFactoryID}

func Test_Users(t *testing.T) {
	assert := assert.New(t)

	ctx := contextWithFactoryIDs(context.Background(), factoryIDs...)

	{ // CreateUsers: insufficient request
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncCreateUsers,
			Input: mock.Input{
				Request: mcom.CreateUsersRequest{
					Users: []mcom.User{},
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

		_, err = mockServer.CreateUsers(ctx, &mesync.CreateUsersRequest{})
		assert.Error(err)
		assert.Equal(parseError(testFactoryID, mcomErr.Error{
			Code: mcomErr.Code_INSUFFICIENT_REQUEST,
		}), err)
	}
	{ // CreateUsers: good case
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncCreateUsers,
			Input: mock.Input{
				Request: mcom.CreateUsersRequest{
					Users: []mcom.User{
						{
							ID:           testUserA,
							DepartmentID: testDepartmentA,
						}, {
							ID:           testUserB,
							DepartmentID: testDepartmentB,
						}, {
							ID:           testADUser,
							Account:      testADAccount,
							DepartmentID: testDepartmentB,
						},
					},
				},
			},
			Output: mock.Output{
				Response: nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.CreateUsers(ctx, &mesync.CreateUsersRequest{
			Users: []*mesync.UserInfo{
				{
					Id:           testUserA,
					DepartmentId: testDepartmentA,
				}, {
					Id:           testUserB,
					DepartmentId: testDepartmentB,
				}, {
					Id:           testADUser,
					Account:      testADAccount,
					DepartmentId: testDepartmentB,
				},
			},
		})
		assert.NoError(err)
	}
	{ // CreateDepartments: insufficient request
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncCreateDepartments,
			Input: mock.Input{
				Request: mcom.CreateDepartmentsRequest{},
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

		_, err = mockServer.CreateDepartments(ctx, &mesync.Departments{
			Ids: []string{},
		})
		assert.Error(err)
		assert.Equal(parseError(testFactoryID, mcomErr.Error{
			Code: mcomErr.Code_INSUFFICIENT_REQUEST,
		}), err)
	}
	{ // CreateDepartments: good case
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncCreateDepartments,
			Input: mock.Input{
				Request: mcom.CreateDepartmentsRequest{testDepartmentA, testDepartmentB},
			},
			Output: mock.Output{
				Response: nil,
				Error:    nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.CreateDepartments(ctx, &mesync.Departments{
			Ids: []string{testDepartmentA, testDepartmentB},
		})
		assert.NoError(err)
	}
	{ // UpdateUser: insufficient request
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncUpdateUser,
			Input: mock.Input{
				Request: mcom.UpdateUserRequest{
					ID:           "",
					DepartmentID: "",
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

		_, err = mockServer.UpdateUser(ctx, &mesync.UpdateUserRequest{
			Id:           "",
			DepartmentId: "",
		})
		assert.Error(err)
		assert.Equal(parseError(testFactoryID, mcomErr.Error{
			Code: mcomErr.Code_INSUFFICIENT_REQUEST,
		}), err)
	}
	{ // UpdateUser: good case for update department
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncUpdateUser,
			Input: mock.Input{
				Request: mcom.UpdateUserRequest{
					ID:           testUserA,
					DepartmentID: testDepartmentA,
				},
			},
			Output: mock.Output{
				Response: nil,
				Error:    nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.UpdateUser(ctx, &mesync.UpdateUserRequest{
			Id:           testUserA,
			DepartmentId: testDepartmentA,
		})
		assert.NoError(err)
	}
	{ // UpdateUser: good case for update leave date
		timeNow := time.Now().Local()
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncUpdateUser,
			Input: mock.Input{
				Request: mcom.UpdateUserRequest{
					ID:        testUserA,
					LeaveDate: timeNow,
				},
			},
			Output: mock.Output{
				Response: nil,
				Error:    nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		ts, err := ptypes.TimestampProto(timeNow)
		assert.NoError(err)
		_, err = mockServer.UpdateUser(ctx, &mesync.UpdateUserRequest{
			Id:        testUserA,
			LeaveDate: ts,
		})
		assert.NoError(err)
	}
	{ // UpdateUser: good case for update user account
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncUpdateUser,
			Input: mock.Input{
				Request: mcom.UpdateUserRequest{
					ID:      testADUser,
					Account: testADNewAccount,
				},
			},
			Output: mock.Output{
				Response: nil,
				Error:    nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		assert.NoError(err)
		_, err = mockServer.UpdateUser(ctx, &mesync.UpdateUserRequest{
			Id:      testADUser,
			Account: testADNewAccount,
		})
		assert.NoError(err)
	}
}
