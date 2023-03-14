package router

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/manager/upstream"
	tcprouter "api_gateway/internal/gateway/router/tcp"
	"api_gateway/pkg/udp"
	"context"
	"crypto/tls"
	"github.com/rs/zerolog/log"
)

// Factory the factory of TCP/UDP routers.
type Factory struct {
	staticConfiguration config.Gateway
	upstreamFactory     *upstream.Factory
	cancelPrevState     context.CancelFunc
}

// NewRouterFactory creates a new RouterFactory.
func NewRouterFactory(staticConfiguration config.Gateway, upstreamFactory *upstream.Factory) *Factory {

	return &Factory{
		staticConfiguration: staticConfiguration,
		upstreamFactory:     upstreamFactory,
	}
}

// CreateTCPRouters creates new TCPRouter.
func (f *Factory) CreateTCPRouters(rootCtx context.Context, rtConf *config.Endpoint) (*tcprouter.Router, *GrpcServer, *GrpcServer) {
	if f.cancelPrevState != nil {
		f.cancelPrevState()
	}
	var (
		httpTLSConfig *tls.Config
		err           error
		ctx           context.Context
	)

	ctx, f.cancelPrevState = context.WithCancel(rootCtx)

	// build http handler
	handlersNonTLS := f.buildHttpHandlers(ctx, rtConf, false)
	handlersTLS := f.buildHttpHandlers(ctx, rtConf, true)

	router, err := tcprouter.NewRouter()
	if err != nil {
		return nil, nil, nil
	}

	if len(getHTTPRouters(rtConf, true)) > 0 {
		httpTLSConfig, err = f.generateHTTPSConfig(rtConf)
		if err != nil {
			log.Debug().Err(err).Msg("Error generate https certs")
		}
	}

	// add http handler to tcp mux
	router.SetHTTPHandler(handlersNonTLS)
	router.SetHTTPSHandler(handlersTLS, httpTLSConfig)

	// build tcp handler
	f.buildTCPHandlers(ctx, router, rtConf)

	// build grpc handler && middleware
	grpcServer := f.buildGrpcHandlers(ctx, rtConf, false)
	grpcTLSServer := f.buildGrpcHandlers(ctx, rtConf, true)

	return router, grpcServer, grpcTLSServer
}

// CreateUDPHandlers creates new UDP Handlers.
func (f *Factory) CreateUDPHandlers(ctx context.Context, rtConf *config.Endpoint) udp.Handler {
	handlers := f.buildUDPHandlers(ctx, rtConf)
	if len(handlers) > 0 {
		return handlers[0]
	}
	return nil
}
