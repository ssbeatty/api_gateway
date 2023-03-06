package config

import "fmt"

const (
	EndpointTypeTCP EndpointType = "tcp"
	EndpointTypeUDP EndpointType = "udp"

	RuleTypeTCP   RuleType = "tcp"
	RuleTypeHTTP  RuleType = "http"
	RuleTypeGRPC  RuleType = "grpc"
	RuleTypeHTTPS RuleType = "https"
	RuleTypeUDP   RuleType = "udp"

	UpstreamTypeURL    UpstreamType = "url"
	UpstreamTypeSTATIC UpstreamType = "static"
	UpstreamTypeServer UpstreamType = "server"
)

// EndpointType
// string: tcp || udp
type EndpointType string

// RuleType
// string: tcp || http || grpc || udp || https
type RuleType string

// UpstreamType
// string: url || static || server
type UpstreamType string

// Endpoint Every Endpoint provide a port
type Endpoint struct {
	Name       string       `yaml:"name"`
	ListenPort int          `yaml:"listen_port"`
	Type       EndpointType `yaml:"type"`
	Routers    []Routers    `yaml:"routers"`
}

func (e *Endpoint) GetAddress() string {
	return fmt.Sprintf("0.0.0.0:%d", e.ListenPort)
}

// Routers host match router
// if not tls enable Host default *
// else got a 4 layer host info use tls
type Routers struct {
	Host        string       `yaml:"host"`
	Rules       []Rule       `yaml:"router"`
	Middlewares []Middleware `yaml:"middlewares"`
	TLSConfig   TLS          `yaml:"tls_config"`
	TlsEnabled  bool         `yaml:"tls_enabled"`
}

// Rule example nginx location
// if rule is http or https can be many
// if rule is grpc or tcp stream can be one
type Rule struct {
	Type     RuleType `yaml:"type"`
	Rule     string   `yaml:"rule"`
	Priority int      `yaml:"priority"`
	Upstream Upstream `yaml:"upstream"`
}

// TLS config
type TLS struct {
	Enable  bool
	CsrPath string
	KeyPath string
}

// Upstream can be file path, url or server with port
type Upstream struct {
	Type  UpstreamType `yaml:"type"`
	Paths []string
}

// Middleware name and config use interface
type Middleware struct {
	Name   string
	Config string
}
