package server

import (
	"context"

	"google.golang.org/grpc/metadata"

	"gitlab.kenda.com.tw/kenda/mcom"
	"gitlab.kenda.com.tw/kenda/mcom/mock"

	grpc_context "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/context"
	grpc_md "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/metadata"
)

const (
	testFactoryID = "TEST"
)

func contextWithFactoryIDs(ctx context.Context, factoryIDs ...string) context.Context {
	md := make(metadata.MD)
	md.Append(grpc_md.FactoryIDs, factoryIDs...)

	parser := grpc_context.NewParser(ctx, md)
	return parser.Parse(grpc_md.FactoryIDs, grpc_context.FactoryIDs, grpc_context.WithMultiple()).Done()
}

func newMockServer(s []mock.Script) (*Server, error) {
	mockDM, err := mock.New(s)
	if err != nil {
		return nil, err
	}
	return &Server{
		dm: map[string]mcom.DataManager{
			testFactoryID: mockDM,
		},
	}, nil
}
