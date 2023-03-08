package backend

import (
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/pkg/safe"
)

func NewBackend() *Backend {
	return &Backend{
		applyMessage: make(chan dynamic.Message, 100),
	}
}

type Backend struct {
	applyMessage chan dynamic.Message
}

func (b *Backend) Provide(configurationChan chan<- dynamic.Message, pool *safe.Pool) error {
	pool.Go(func() {
		for msg := range b.applyMessage {
			configurationChan <- msg
		}
	})
	return nil
}

func (b *Backend) ReloadConfig(msg dynamic.Message) {
	b.applyMessage <- msg
}

func (b *Backend) Init() error {
	return nil
}

func (b *Backend) Name() string {
	return "backend"
}
