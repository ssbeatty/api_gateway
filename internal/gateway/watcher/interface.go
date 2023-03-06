package watcher

import (
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/pkg/safe"
)

type Provider interface {
	// Provide allows the provider to provide configurations to gateway
	// using the given configuration channel.
	Provide(configurationChan chan<- dynamic.Message, pool *safe.Pool) error
	Name() string
	Init() error
}
