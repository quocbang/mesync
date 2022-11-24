package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/xid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	commonsCtx "gitlab.kenda.com.tw/kenda/commons/v2/utils/context"
	"gitlab.kenda.com.tw/kenda/mcom"
	"gitlab.kenda.com.tw/kenda/mcom/impl"

	"gitlab.kenda.com.tw/kenda/mesync/pkg/config"
	pb "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
	grpc_context "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/context"
	"gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/interceptors"
	"gitlab.kenda.com.tw/kenda/mesync/server"
)

type Server struct {
	Host string
	Port int
	TLS  config.TLSOptionsType

	// Config defines data storage endpoint.
	// It is required to set.
	Config config.Config

	// Timeout is each request timeout.
	Timeout time.Duration

	// AuthKey is authentication key to authenticate the request.
	AuthKey string
}

// Run serves gRPC server.
func (s *Server) Run() error {
	defer zap.L().Sync() // nolint

	opts, err := interceptors.Config{
		GenerateRequestID: generateRequestID,
		ParseLogFields:    getContextFields,
		WithLoggerFunc:    contextWithLogger,
		Timeout:           s.Timeout,
		AuthKey:           s.AuthKey,
	}.Build()
	if err != nil {
		return err
	}

	if s.TLS.UseTLS() {
		cred, err := credentials.NewServerTLSFromFile(s.TLS.Cert, s.TLS.Key)
		if err != nil {
			return err
		}
		opts = append(opts, grpc.Creds(cred))
	}

	mesServer := grpc.NewServer(opts...)

	dms, cloudDM, err := buildDataManagers(s.Config)
	if err != nil {
		return err
	}

	if s.Config.Kenda != nil {
		activeServer := server.NewServer(server.ServerParams{
			LocalDataManager: dms,
			Logger: func(ctx context.Context) *zap.Logger {
				return commonsCtx.Logger(ctx)
			},
		})

		pb.RegisterMesyncServer(mesServer, activeServer)
	}

	if s.Config.Cloud != nil {
		activeCloudServer := server.NewCloudServer(server.CloudServerParams{
			CloudDataManager: cloudDM,
			CloudConfigs:     *s.Config.Cloud,
			Logger: func(ctx context.Context) *zap.Logger {
				return commonsCtx.Logger(ctx)
			},
		})

		activeCloudServer.StartReUploadTicker()

		pb.RegisterCloudServer(mesServer, activeCloudServer)
	}

	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	zap.L().Info("serving gRPC server..", zap.String("address", address))

	return mesServer.Serve(lis)
}

func contextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return commonsCtx.WithLogger(ctx, logger)
}

func generateRequestID() string {
	return xid.New().String()
}

func buildDataManagers(conf config.Config) (dms map[string]mcom.DataManager, cloudDM mcom.DataManager, err error) {
	if conf.Kenda != nil {
		dms = make(map[string]mcom.DataManager, len(*conf.Kenda))
		for factory, dbConfig := range *conf.Kenda {
			dm, err := getDM(dbConfig)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to build DataManager for %s: %v", factory, err)
			}

			dms[factory] = dm
		}
	}

	if conf.Cloud != nil {
		cloudDM, err = buildCloudDataManager(*conf.Cloud)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build DataManager for cloud service: %v", err)
		}
	}

	return
}

func buildCloudDataManager(conf config.CloudConfigs) (mcom.DataManager, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	return getDM(conf.Azure)
}

func getDM(dbConfig config.Database) (mcom.DataManager, error) {
	opt := impl.WithPostgreSQLSchema(dbConfig.Schema)
	pgCfg := impl.PGConfig{
		Database: dbConfig.Name,
		Address:  dbConfig.Address,
		Port:     dbConfig.Port,
		UserName: dbConfig.UserName,
		Password: dbConfig.Password,
	}

	return impl.New(context.Background(), pgCfg, opt)
}

// getContextFields return corresponding useful header request for logging records.
func getContextFields(ctx context.Context) []zap.Field {
	return []zap.Field{
		zap.String("request_id", grpc_context.GetRequestID(ctx)),
		zap.String("user_id", grpc_context.GetUserID(ctx)),
		zap.Strings("factory_ids", grpc_context.GetFactoryIDs(ctx)),
	}
}
