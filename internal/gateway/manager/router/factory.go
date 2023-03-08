package router

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/manager/upstream"
	tcprouter "api_gateway/internal/gateway/router/tcp"
	"api_gateway/pkg/udp"
	"context"
)

// Factory the factory of TCP/UDP routers.
type Factory struct {
	staticConfiguration config.Gateway
	upstreamFactory     *upstream.Factory
}

// NewRouterFactory creates a new RouterFactory.
func NewRouterFactory(staticConfiguration config.Gateway, upstreamFactory *upstream.Factory) *Factory {

	return &Factory{
		staticConfiguration: staticConfiguration,
		upstreamFactory:     upstreamFactory,
	}
}

// CreateTCPRouters creates new TCPRouter.
func (f *Factory) CreateTCPRouters(ctx context.Context, rtConf *config.Endpoint) (*tcprouter.Router, *GrpcServer) {

	// build http handler
	handlersNonTLS := f.buildHttpHandlers(ctx, rtConf, false)
	handlersTLS := f.buildHttpHandlers(ctx, rtConf, true)

	router, err := tcprouter.NewRouter()
	if err != nil {
		return nil, nil
	}

	// add http handler to tcp mux
	router.SetHTTPHandler(handlersNonTLS)
	router.SetHTTPSHandler(handlersTLS, rtConf.TLSConfig.Config)

	// build tcp handler
	err = f.buildTCPHandlers(ctx, router, rtConf)
	if err != nil {
		return nil, nil
	}

	// build grpc handler && middleware
	grpcServer := f.buildGrpcHandlers(ctx, rtConf)

	return router, grpcServer
}

// CreateUDPHandlers creates new UDP Handlers.
func (f *Factory) CreateUDPHandlers(ctx context.Context, rtConf *config.Endpoint) udp.Handler {
	handlers := f.buildUDPHandlers(ctx, rtConf)
	if len(handlers) > 0 {
		return handlers[0]
	}
	return nil
}
