package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr"`
	Token  string `ocisConfig:"token"`
	Pprof  bool   `ocisConfig:"pprof"`
	Zpages bool   `ocisConfig:"zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr"`
	Namespace string `ocisConfig:"namespace"`
	Root      string `ocisConfig:"root"`
}

// Server configures a server.
type Server struct {
	Version string `ocisConfig:"version"`
	Name    string `ocisConfig:"name"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `ocisConfig:"address"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret"`
}

type Spaces struct {
	WebDavBase   string `ocisConfig:"webdav_base"`
	WebDavPath   string `ocisConfig:"webdav_path"`
	DefaultQuota string `ocisConfig:"default_quota"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File         string       `ocisConfig:"file"`
	Log          *shared.Log  `ocisConfig:"log"`
	Debug        Debug        `ocisConfig:"debug"`
	HTTP         HTTP         `ocisConfig:"http"`
	Server       Server       `ocisConfig:"server"`
	Tracing      Tracing      `ocisConfig:"tracing"`
	Reva         Reva         `ocisConfig:"reva"`
	TokenManager TokenManager `ocisConfig:"token_manager"`
	Spaces       Spaces       `ocisConfig:"spaces"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:  "127.0.0.1:9124",
			Token: "",
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9120",
			Namespace: "com.owncloud.web",
			Root:      "/graph",
		},
		Server: Server{},
		Tracing: Tracing{
			Enabled: false,
			Type:    "jaeger",
			Service: "graph",
		},
		Reva: Reva{
			Address: "127.0.0.1:9142",
		},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		Spaces: Spaces{
			WebDavBase:   "https://localhost:9200",
			WebDavPath:   "/dav/spaces/",
			DefaultQuota: "1000000000",
		},
	}
}
