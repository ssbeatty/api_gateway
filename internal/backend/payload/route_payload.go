package payload

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
)

// request

type PostEndPointReq struct {
	Name    string              `json:"name" binding:"required" example:"tcp_endpoint_1"`
	Type    config.EndpointType `json:"type" binding:"required" example:"tcp"`
	Routers []RouterInfo        `json:"routers"`
}

type RouterInfo struct {
	Rule        string           `json:"rule" binding:"required" example:"Host(\"api.demo.com\") && PathPrefix(\"/\")"`
	Type        config.RuleType  `json:"type" binding:"required" example:"http"`
	TlsEnable   bool             `json:"tls_enable"`
	Priority    int              `json:"priority"`
	Host        string           `json:"host" binding:"required_if=TlsEnable true,hostname"`
	UpStream    UpStreamInfo     `json:"up_stream" binding:"required"`
	Tls         TlsInfo          `json:"tls"  binding:"required_if=TlsEnable true"`
	Middlewares []MiddleWareInfo `json:"middlewares"`
}

type MiddleWareInfo struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
}

type UpStreamInfo struct {
	Type        config.UpstreamType `json:"type" binding:"required" example:"url"`
	Path        []string            `json:"path"`
	Weights     []string            `json:"weights"`
	LoadBalance loadbalancer.LbType `json:"load_balance"`
}

type TlsInfo struct {
	Type       string `json:"type"`
	ClientAuth string `json:"client_auth"`
}

type Pages struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}
