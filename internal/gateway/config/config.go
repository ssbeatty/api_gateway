package config

import (
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
	"fmt"
)

const (
	EndpointTypeTCP EndpointType = "tcp"
	EndpointTypeUDP EndpointType = "udp"

	RuleTypeTCP  RuleType = "tcp"
	RuleTypeHTTP RuleType = "http"
	RuleTypeGRPC RuleType = "grpc"
	RuleTypeUDP  RuleType = "udp"

	UpstreamTypeURL    UpstreamType = "url"
	UpstreamTypeSTATIC UpstreamType = "static"
	UpstreamTypeServer UpstreamType = "server"

	TLSTypeBytes TLSType = "bytes"
	TLSTypePath  TLSType = "file"
)

// EndpointType
// string: tcp || udp
type EndpointType string

// RuleType
// string: tcp || http || grpc || udp || https
type RuleType string

// UpstreamType
// string: url || static || server || service configs
type UpstreamType string

// TLSType bytes || path
type TLSType string

// Endpoint Every Endpoint provide a port
type Endpoint struct {
	Name       string       `yaml:"name" mapstructure:"name"`
	ListenPort int          `yaml:"listen_port" mapstructure:"listen_port"`
	Type       EndpointType `yaml:"type" mapstructure:"type"`
	Routers    []Router     `yaml:"routers" mapstructure:"routers"`
}

func (e *Endpoint) GetAddress() string {
	return fmt.Sprintf("0.0.0.0:%d", e.ListenPort)
}

// Router host match router
// if not tls enable Host default *
// else got a 4 layer host info use tls
type Router struct {
	Rule        string       `yaml:"rule" mapstructure:"rule"`
	Host        string       `yaml:"host" mapstructure:"host"`
	TlsEnabled  bool         `yaml:"tls_enabled" mapstructure:"tls_enabled"`
	Type        RuleType     `yaml:"type" mapstructure:"type"`
	Priority    int          `yaml:"priority" mapstructure:"priority"`
	Middlewares []Middleware `yaml:"middlewares" mapstructure:"middlewares"`
	Upstream    Upstream     `yaml:"upstream" mapstructure:"upstream"`
	// every host has one tls config
	TLSConfig TLS `yaml:"tls_config" mapstructure:"tls_config"`
}

// TLS config
type TLS struct {
	Type       TLSType  `yaml:"type" mapstructure:"type"`
	CsrFile    string   `yaml:"csr_file" mapstructure:"csr_file"`
	KeyFile    string   `yaml:"key_file" mapstructure:"key_file"`
	CaFiles    []string `yaml:"ca_files" mapstructure:"ca_files"`
	ClientAuth string   `yaml:"client_auth" mapstructure:"client_auth"`
}

// Upstream can be file path, url or server with port
type Upstream struct {
	Type             UpstreamType        `yaml:"type" mapstructure:"type"`
	Paths            []string            `yaml:"paths" mapstructure:"paths"`
	Weights          []int               `yaml:"weights" mapstructure:"weights"`
	LoadBalancerType loadbalancer.LbType `yaml:"load_balance" mapstructure:"load_balance"`
}

// Middleware name and config use interface
// 4 layer middleware example tcp
// 7 layer middleware example http or https or grpc
type Middleware struct {
	Name   string                 `yaml:"name" mapstructure:"name"`
	Type   string                 `yaml:"type" mapstructure:"type"`
	Config map[string]interface{} `yaml:"config" mapstructure:"config"`
}
