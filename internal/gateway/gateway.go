package gateway

import (
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/internal/gateway/watcher"
	"api_gateway/pkg/safe"
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type Server struct {
	watcher      *watcher.ConfigurationWatcher
	routinesPool *safe.Pool

	stopChan chan struct{}
}

func NewServer(routinesPool *safe.Pool, watcher *watcher.ConfigurationWatcher) *Server {

	return &Server{
		watcher:      watcher,
		routinesPool: routinesPool,
		stopChan:     make(chan struct{}, 1),
	}
}

// Start starts the server and Stop/Close it when context is Done.
func (s *Server) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		logger := log.Ctx(ctx)
		logger.Info().Msg("I have to go...")
		logger.Info().Msg("Stopping server gracefully")
		s.Stop()
	}()

	s.setupConfigWatcher()
	s.watcher.Start()
}

// Wait blocks until the server shutdown.
func (s *Server) Wait() {
	<-s.stopChan
}

// Stop stops the server.
func (s *Server) Stop() {
	defer log.Info().Msg("Server stopped")

	s.stopChan <- struct{}{}
}

// Close destroys the server.
func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	go func(ctx context.Context) {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			panic("Timeout while stopping, killing instance ✝")
		}
	}(ctx)

	s.routinesPool.Stop()

	cancel()
}

func (s *Server) setupConfigWatcher() {
	s.watcher.AddListener(func(configuration dynamic.Configuration) {
		switch configuration.Action {
		case watcher.ActionCreate:

		}
	})
}
