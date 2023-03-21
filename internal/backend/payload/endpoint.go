package payload

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
)

// PostEndPointReq for create endpoint
type PostEndPointReq struct {
	Name       string              `json:"name" binding:"required"`
	Type       config.EndpointType `json:"type" binding:"required,oneof=tcp udp"`
	ListenPort int                 `json:"listen_port"`
	Routers    []RouterInfo        `json:"routers"`
}

// OptionEndpointReq for option endpoint, uri id
type OptionEndpointReq struct {
	Id int `uri:"id" binding:"required"`
}

// RouterInfo router payload
type RouterInfo struct {
	Id          int              `json:"id"`
	Rule        string           `json:"rule" binding:"required"`
	Type        config.RuleType  `json:"type" binding:"required,oneof=tcp udp http grpc"`
	TlsEnable   bool             `json:"tls_enable"`
	Priority    int              `json:"priority"`
	Host        string           `json:"host" binding:"required_if=TlsEnable true,hostname"`
	UpStream    UpstreamInfo     `json:"upstream" binding:"required"`
	CertId      int              `json:"cert_id"  binding:"required_if=TlsEnable true"`
	Middlewares []MiddleWareInfo `json:"middlewares"`
}

// MiddleWareInfo middleware payload
type MiddleWareInfo struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
}

// UpstreamInfo upstream payload
type UpstreamInfo struct {
	Type                config.UpstreamType `json:"type" binding:"required,oneof=url static server"`
	Path                []string            `json:"path"`
	Weights             []int               `json:"weights"`
	LoadBalance         loadbalancer.LbType `json:"load_balance"`
	MaxIdleConnsPerHost int                 `json:"maxIdleConnsPerHost,omitempty"`
}
