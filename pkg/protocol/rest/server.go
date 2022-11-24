package rest

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"gitlab.kenda.com.tw/kenda/mesync/pkg/config"
	pb "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
)

type Server struct {
	Host string
	Port int
	TLS  config.TLSOptionsType

	GRPCServerEndpoint string

	// EnableMESService toggles if serve the service to synchronize
	// data with MES.
	EnableMESService bool

	// EnableCloudService toggles if serve the service to communicate
	// with the cloud.
	EnableCloudService bool
}

func (s *Server) gatewayMux() (*runtime.ServeMux, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(h string) (string, bool) {
			return h, true
		}),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
	)

	var opts []grpc.DialOption

	if s.TLS.UseTLS() {
		var cred credentials.TransportCredentials
		if s.TLS.SpecifiedRootCA() {
			var err error
			cred, err = credentials.NewClientTLSFromFile(s.TLS.RootCA, "")
			if err != nil {
				return nil, cancel, err
			}
		} else {
			cred = credentials.NewTLS(&tls.Config{})
		}
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(cred),
		}
	} else {
		opts = []grpc.DialOption{
			grpc.WithInsecure(),
		}
	}

	// Register MESync handlers
	if s.EnableMESService {
		if err := pb.RegisterMesyncHandlerFromEndpoint(ctx, mux, s.GRPCServerEndpoint, opts); err != nil {
			return nil, cancel, err
		}
	}

	// Register Cloud handlers
	if s.EnableCloudService {
		if err := pb.RegisterCloudHandlerFromEndpoint(ctx, mux, s.GRPCServerEndpoint, opts); err != nil {
			return nil, cancel, err
		}
	}

	return mux, cancel, nil
}

// Run serves RESTful server.
func (s *Server) Run() error {
	defer zap.L().Sync() // nolint

	mux, cancel, err := s.gatewayMux()
	if err != nil {
		return err
	}
	defer cancel()

	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	zap.L().Info("serving gateway server..", zap.String("address", address))

	if s.TLS.UseTLS() {
		return srv.ListenAndServeTLS(s.TLS.Cert, s.TLS.Key)
	}

	return srv.ListenAndServe()
}
