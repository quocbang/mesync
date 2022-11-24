package server

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"

	commonsCtx "gitlab.kenda.com.tw/kenda/commons/v2/utils/context"
	"gitlab.kenda.com.tw/kenda/mcom"

	"gitlab.kenda.com.tw/kenda/mesync/pkg/config"
	grpc_context "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/context"
)

type eachFactoryFunc func(ctx context.Context, factoryID string, i int, dm mcom.DataManager) error

// Server definition.
type Server struct {
	dm        map[string]mcom.DataManager
	getLogger func(context.Context) *zap.Logger
}

type CloudServer struct {
	cloud     cloud
	getLogger func(context.Context) *zap.Logger
}

type cloud struct {
	dm      mcom.DataManager
	configs config.CloudConfigs
}

// NewServer create new server.
func NewServer(params ServerParams) Server {
	return Server{
		dm:        params.LocalDataManager,
		getLogger: params.Logger,
	}
}

func NewCloudServer(params CloudServerParams) CloudServer {
	return CloudServer{
		cloud: cloud{
			dm:      params.CloudDataManager,
			configs: params.CloudConfigs,
		},
		getLogger: params.Logger,
	}
}

type ServerParams struct {
	LocalDataManager map[string]mcom.DataManager
	Logger           func(context.Context) *zap.Logger
}

type CloudServerParams struct {
	CloudDataManager mcom.DataManager
	CloudConfigs     config.CloudConfigs
	Logger           func(context.Context) *zap.Logger
}

func (s Server) getDataManagerByFactory(id string) (mcom.DataManager, error) {
	dm, ok := s.dm[id]
	if !ok {
		return nil, fmt.Errorf("mcom service not found [%s]", id)
	}
	return dm, nil
}

// CheckServer checks if server is active.
func (s Server) CheckServer(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// eachFactory runs the f function for each factory from the context.
// eachFactory runs each factory concurrently.
func (s Server) eachFactory(ctx context.Context, f eachFactoryFunc) error {
	ids := grpc_context.GetFactoryIDs(ctx)
	if len(ids) == 0 {
		return fmt.Errorf("missing factory id")
	}

	ctx = commonsCtx.WithUserID(ctx, grpc_context.GetUserID(ctx))

	ch := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(len(ids))
	for i, id := range ids {
		go func(i int, id string) {
			defer wg.Done()

			dm, err := s.getDataManagerByFactory(id)
			if err != nil {
				ch <- parseError(id, err)
				return
			}

			if err := f(ctx, id, i, dm); err != nil {
				ch <- parseError(id, err)
			}
		}(i, id)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	return mergeErrors(ch)
}

func parseError(name string, err error) error {
	return fmt.Errorf("%v [%s]", err, name)
}

func mergeErrors(ch <-chan error) error {
	ss := []string{}
	for err := range ch {
		ss = append(ss, err.Error())
	}

	if len(ss) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(ss, "; "))
}
