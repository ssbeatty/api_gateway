package router

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/manager/upstream"
	tcprouter "api_gateway/internal/gateway/router/tcp"
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
func (f *Factory) CreateTCPRouters(ctx context.Context, rtConf *config.Endpoint) *tcprouter.Router {

	handlersNonTLS := f.buildHttpHandlers(ctx, rtConf, false)
	handlersTLS := f.buildHttpHandlers(ctx, rtConf, true)

	router, err := tcprouter.NewRouter()
	if err != nil {
		return nil
	}

	router.SetHTTPHandler(handlersNonTLS)
	router.SetHTTPSHandler(handlersTLS, rtConf.TLSConfig.Config)

	_, err = f.buildTCPHandlers(ctx, router, rtConf)
	if err != nil {
		return nil
	}
	return router
}
