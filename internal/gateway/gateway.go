package gateway

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/internal/gateway/endpoint/tcp"
	"api_gateway/internal/gateway/endpoint/udp"
	routerManager "api_gateway/internal/gateway/manager/router"
	"api_gateway/internal/gateway/watcher"
	"api_gateway/pkg/safe"
	"context"
	"github.com/rs/zerolog/log"
)

// Server Gateway main server
type Server struct {
	watcher       *watcher.ConfigurationWatcher
	routinesPool  *safe.Pool
	tcpEndpoints  *safe.SyncMap[string, *tcp.EndPoint]
	udpEndpoints  *safe.SyncMap[string, *udp.EndPoint]
	stopChan      chan struct{}
	routerFactory *routerManager.Factory
	gatewayConfig *config.Gateway
}

func NewServer(
	routinesPool *safe.Pool,
	watcher *watcher.ConfigurationWatcher,
	routerFactory *routerManager.Factory,
	gatewayConfig *config.Gateway,
) *Server {

	return &Server{
		watcher:       watcher,
		routinesPool:  routinesPool,
		stopChan:      make(chan struct{}, 1),
		tcpEndpoints:  safe.NewSyncMap[string, *tcp.EndPoint](),
		udpEndpoints:  safe.NewSyncMap[string, *udp.EndPoint](),
		routerFactory: routerFactory,
		gatewayConfig: gatewayConfig,
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

	s.tcpEndpoints.Range(func(_ string, endpoint *tcp.EndPoint) bool {
		endpoint.Shutdown(ctx)

		return true
	})
	s.udpEndpoints.Range(func(_ string, endpoint *udp.EndPoint) bool {
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
					log.Error().Msgf("Error when New TCP Endpoint, %v", err)
				}
				s.routinesPool.GoCtx(endpoint.Start)
				s.switchTCPRouter(ctx, endpoint, configuration)

				s.tcpEndpoints.Store(endpointConfig.Name, endpoint)
			} else if endpointConfig.Type == config.EndpointTypeUDP {
				endpoint, err := udp.NewUDPEntryPoint(&configuration.EndPoint, s.gatewayConfig, s.routinesPool)
				if err != nil {
					log.Error().Msgf("Error when New UDP Endpoint, %v", err)
				}
				s.routinesPool.GoCtx(endpoint.Start)
				s.switchUDPRouter(ctx, endpoint, configuration)

				s.udpEndpoints.Store(endpointConfig.Name, endpoint)
			}

		case watcher.ActionDelete:
			if endpoint, ok := s.tcpEndpoints.Load(endpointConfig.Name); ok {
				endpoint.Shutdown(ctx)

				s.tcpEndpoints.Delete(endpointConfig.Name)
			}

			if endpoint, ok := s.udpEndpoints.Load(endpointConfig.Name); ok {
				endpoint.Shutdown(ctx)

				s.udpEndpoints.Delete(endpointConfig.Name)
			}
		case watcher.ActionUpdate:
			if endpoint, ok := s.tcpEndpoints.Load(endpointConfig.Name); ok {
				s.switchTCPRouter(ctx, endpoint, configuration)
			}
			if endpoint, ok := s.udpEndpoints.Load(endpointConfig.Name); ok {
				s.switchUDPRouter(ctx, endpoint, configuration)
			}
		}
	})
}

func (s *Server) switchTCPRouter(ctx context.Context, serverEntryPointsTCP *tcp.EndPoint, conf dynamic.Configuration) {
	routers, grpcServer := s.routerFactory.CreateTCPRouters(ctx, &conf.EndPoint)

	serverEntryPointsTCP.SwitchRouter(routers, grpcServer)
}

func (s *Server) switchUDPRouter(ctx context.Context, serverEntryPointsUDP *udp.EndPoint, conf dynamic.Configuration) {
	handlers := s.routerFactory.CreateUDPHandlers(ctx, &conf.EndPoint)

	serverEntryPointsUDP.Switch(handlers)
}
