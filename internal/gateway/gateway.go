package gateway

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/internal/gateway/endpoint/tcp"
	routerManager "api_gateway/internal/gateway/manager/router"
	tcprouter "api_gateway/internal/gateway/router/tcp"
	"api_gateway/internal/gateway/watcher"
	"api_gateway/pkg/safe"
	"context"
	"github.com/rs/zerolog/log"
)

type Endpoint interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	SwitchRouter(rt *tcprouter.Router)
}

type Server struct {
	watcher       *watcher.ConfigurationWatcher
	routinesPool  *safe.Pool
	endpoints     *safe.SyncMap[string, Endpoint]
	stopChan      chan struct{}
	routerFactory *routerManager.Factory
}

func NewServer(routinesPool *safe.Pool, watcher *watcher.ConfigurationWatcher, routerFactory *routerManager.Factory) *Server {

	return &Server{
		watcher:       watcher,
		routinesPool:  routinesPool,
		stopChan:      make(chan struct{}, 1),
		endpoints:     safe.NewSyncMap[string, Endpoint](),
		routerFactory: routerFactory,
	}
}

// Start starts the server and Stop/Close it when context is Done.
func (s *Server) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		logger := log.Ctx(ctx)
		logger.Info().Msg("I have to go...")
		logger.Info().Msg("Stopping server gracefully")
		s.Stop(ctx)
	}()

	s.setupConfigWatcher(ctx)
	s.watcher.Start()
}

// Wait blocks until the server shutdown.
func (s *Server) Wait() {
	<-s.stopChan
}

// Stop stops the server.
func (s *Server) Stop(ctx context.Context) {
	defer log.Info().Msg("Server stopped")

	s.endpoints.Range(func(_ string, endpoint Endpoint) bool {
		endpoint.Shutdown(ctx)

		return true
	})
	s.stopChan <- struct{}{}
}

// Close destroys the server.
func (s *Server) Close() {
	s.routinesPool.Stop()
}

func (s *Server) setupConfigWatcher(ctx context.Context) {
	// watch the endpoint action
	s.watcher.AddListener(func(configuration dynamic.Configuration) {
		endpointConfig := configuration.EndPoint
		switch configuration.Action {
		case watcher.ActionCreate:
			if endpointConfig.Type == config.EndpointTypeTCP {
				endpoint, err := tcp.NewTCPEndPoint(ctx, &configuration.EndPoint, s.routinesPool)
				if err != nil {
					log.Error().Msgf("Error when New Tcp Endpoint, %v", err)
				}
				s.routinesPool.GoCtx(endpoint.Start)
				s.switchRouter(ctx, endpoint, configuration)

				s.endpoints.Store(endpointConfig.Name, endpoint)
			}

		case watcher.ActionDelete:
			if endpoint, ok := s.endpoints.Load(endpointConfig.Name); ok {
				endpoint.Shutdown(ctx)

				s.endpoints.Delete(endpointConfig.Name)
			}
		case watcher.ActionUpdate:
			if endpoint, ok := s.endpoints.Load(endpointConfig.Name); ok {
				s.switchRouter(ctx, endpoint, configuration)
			}
		}
	})
}

func (s *Server) switchRouter(ctx context.Context, serverEntryPointsTCP Endpoint, conf dynamic.Configuration) {
	routers := s.routerFactory.CreateTCPRouters(ctx, &conf.EndPoint)

	serverEntryPointsTCP.SwitchRouter(routers)
}
