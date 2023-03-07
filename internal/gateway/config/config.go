package config

import (
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
	middleware "api_gateway/pkg/middlewares"
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
	Name       string       `yaml:"name"`
	ListenPort int          `yaml:"listen_port"`
	Type       EndpointType `yaml:"type"`
	Routers    []Router     `yaml:"routers"`
	TLSConfig  TLS          `yaml:"tls_config"`
}

func (e *Endpoint) GetAddress() string {
	return fmt.Sprintf("0.0.0.0:%d", e.ListenPort)
}

// Router host match router
// if not tls enable Host default *
// else got a 4 layer host info use tls
type Router struct {
	Rule        string       `yaml:"rule"`
	Host        string       `yaml:"host"`
	TlsEnabled  bool         `yaml:"tls_enabled"`
	Type        RuleType     `yaml:"type"`
	Priority    int          `yaml:"priority"`
	Middlewares []Middleware `yaml:"middlewares"`
	Upstream    Upstream     `yaml:"upstream"`
}

// TLS config
type TLS struct {
	Config *tls.Config
}

// Upstream can be file path, url or server with port
type Upstream struct {
	Type             UpstreamType `yaml:"type"`
	Paths            []string
	Weights          []int
	LoadBalancerType loadbalancer.LbType
}

// Middleware name and config use interface
// 4 layer middleware example tcp or udp
// 7 layer middleware example http or https or grpc
type Middleware struct {
	Name   string
	Type   string
	Config middleware.MiddlewareCfg
}
