package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-openapi/swag"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"gitlab.kenda.com.tw/kenda/mesync/pkg/config"
	"gitlab.kenda.com.tw/kenda/mesync/pkg/logger"
	"gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc"
	"gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/rest"
)

// RunServer run gRPC server and RESTful gateway
func RunServer() {
	flags := parseFlags()

	logger.Init(flags.Options.DevMode)

	conf, err := loadConfig(flags.Options.ConfigPath)
	if err != nil {
		zap.L().Fatal("failed to load config", zap.Error(err))
	}

	if !flags.TLSOptions.UseTLS() {
		zap.L().Warn("serving service without TLS handshake")
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()

		s := grpc.Server{
			Host:    flags.Options.GRPCHost,
			Port:    flags.Options.GRPCPort,
			TLS:     flags.TLSOptions,
			Config:  conf,
			Timeout: flags.Options.Timeout,
			AuthKey: flags.Options.AuthKey,
		}

		if err := s.Run(); err != nil {
			zap.L().Fatal("failed to serve gRPC: ", zap.Error(err))
		}
		zap.L().Info("gRPC server stopped")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		s := rest.Server{
			Host:               flags.Options.GatewayHost,
			Port:               flags.Options.GatewayPort,
			TLS:                flags.TLSOptions,
			GRPCServerEndpoint: fmt.Sprintf("%s:%d", flags.Options.GRPCHost, flags.Options.GRPCPort),
			EnableMESService:   conf.Kenda != nil,
			EnableCloudService: conf.Cloud != nil,
		}

		if err := s.Run(); err != nil {
			zap.L().Fatal("failed to serve gRPC gateway: ", zap.Error(err))
		}
		zap.L().Info("gRPC gateway stopped")
	}()
}

type Config struct {
	Options    config.Options
	TLSOptions config.TLSOptionsType
}

// parseFlags parse server configurations.
func parseFlags() *Config {
	var conf Config

	// 1. set available server configurations.
	configurations := []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "Server Configuration",
			LongDescription:  "Server Configuration",
			Options:          &conf.Options,
		},
		{
			ShortDescription: "TLS handshake Configuration",
			LongDescription:  "TLS handshake Configuration",
			Options:          &conf.TLSOptions,
		},
		// more configuration options..
	}

	// 2. Parse command line flags
	parser := flags.NewParser(nil, flags.Default)
	for _, optsGroup := range configurations {
		if _, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options); err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			code = 0
		}
		os.Exit(code)
	}

	return &conf
}

func loadConfig(path string) (config.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return config.Config{}, err
	}
	defer f.Close()

	var conf config.Config
	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		return config.Config{}, err
	}
	return conf, nil
}
