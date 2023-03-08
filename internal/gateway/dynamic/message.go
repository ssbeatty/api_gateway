package dynamic

import "api_gateway/internal/gateway/config"

// Message holds configuration information exchanged between parts of gateway.
type Message struct {
	ProviderName  string
	Configuration []config.Endpoint
}

type Configuration struct {
	Action   string
	EndPoint config.Endpoint
}
