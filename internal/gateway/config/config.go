package config

import "fmt"

// EndpointType
// string: tcp || udp
type EndpointType string

// RuleType
// string: tcp || http || grpc || udp
type RuleType string

// UpstreamType
// string: url || static || server
type UpstreamType string

// Endpoint Every Endpoint provide a port
type Endpoint struct {
	Name       string       `yaml:"name"`
	ListenPort int          `yaml:"listen_port"`
	Type       EndpointType `yaml:"type"`
	Rules      []Rule       `yaml:"rules"`
}

func (e *Endpoint) GetAddress() string {
	return fmt.Sprintf("0.0.0.0:%d", e.ListenPort)
}

// Rule location rule
type Rule struct {
	Type        RuleType     `yaml:"type"`
	Router      Router       `yaml:"router"`
	Upstream    Upstream     `yaml:"upstream"`
	Middlewares []Middleware `yaml:"middlewares"`
	TLS         TLS
}

// Router match router
type Router struct {
	Rule     string `yaml:"rule"`
	Priority int    `yaml:"priority"`
}

// TLS config
type TLS struct {
	Enable  bool
	CsrPath string
	KeyPath string
}

// Upstream upstream
type Upstream struct {
	Type  UpstreamType `yaml:"type"`
	Paths []string
}

// Middleware name and config
type Middleware struct {
	Name   string
	Config string
}
