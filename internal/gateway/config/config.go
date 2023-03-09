package config

import (
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
	"crypto/tls"
	"fmt"
)

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
// string: url || static || server || service configs
type UpstreamType string

// Endpoint Every Endpoint provide a port
type Endpoint struct {
	Name       string       `yaml:"name" mapstructure:"name"`
	ListenPort int          `yaml:"listen_port" mapstructure:"listen_port"`
	Type       EndpointType `yaml:"type" mapstructure:"type"`
	Routers    []Router     `yaml:"routers" mapstructure:"routers"`
	TLSConfig  TLS          `yaml:"tls_config" mapstructure:"tls_config"`
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
}

// TLS config
type TLS struct {
	Config *tls.Config
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
