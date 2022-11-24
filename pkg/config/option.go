package config

import "time"

// Options are configuration flags for Server
type Options struct {
	// gRPC server port, this setting is allowed to avoid situation when the default port is fully bind
	GRPCHost string `long:"grpc-host" description:"the IP to listen on" default:"localhost" env:"MESYNC_GRPC_HOST"`
	GRPCPort int    `long:"grpc-port" description:"this is an internal setting, you could use it when the default port is full on your system" default:"8081" env:"MESYNC_GRPC_PORT"`

	// gRPC Gateway address
	GatewayHost string `long:"gateway-host" description:"the IP to listen on" default:"localhost" env:"MESYNC_GATEWAY_HOST"`
	GatewayPort int    `long:"gateway-port" description:"gateway port to bind" default:"9091" env:"MESYNC_GATEWAY_PORT"`

	// development logging mode
	DevMode bool `long:"dev-mode" description:"development mode" env:"MESYNC_DEV_MODE"`

	// server settings
	// authentication key to send request to api
	AuthKey string `long:"auth-key" description:"authentication key" env:"MESYNC_AUTH_KEY"`
	// each request timeout, default value is 1 minute
	Timeout time.Duration `long:"timeout" description:"request timeout, need to fetch unit time" default:"1m" env:"MESYNC_TIMEOUT"`

	// server settings
	ConfigPath string `long:"config-path" description:"server settings path" env:"MESYNC_CONFIG_PATH"`
}

// TLSOptionsType for enabled TLS handshake configurations.
// In Order to run server using TLS handshake, both Cert and Key must not empty.
type TLSOptionsType struct {
	Cert   string `long:"tls-cert" description:"path to TLS certificate (PUBLIC). To enable TLS handshake, you must set this value" env:"MESYNC_TLS_CERT"`
	Key    string `long:"tls-key" description:"path to TLS certificate key (PRIVATE), To enable TLS handshake, you must set this value" env:"MESYNC_TLS_KEY"`
	RootCA string `long:"root-ca" description:"path to the root certificate"`
}

// UseTLS returns true if there's TLS cert and key assigned.
// In Order to run server using TLS handshake, both `Cert` and `Key` must not be emptied.
func (tls TLSOptionsType) UseTLS() bool {
	return tls.Cert != "" && tls.Key != ""
}

func (tls TLSOptionsType) SpecifiedRootCA() bool {
	return tls.RootCA != ""
}
