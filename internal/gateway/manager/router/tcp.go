package router

import (
	"api_gateway/internal/gateway/config"
	tcprouter "api_gateway/internal/gateway/router/tcp"
	"api_gateway/pkg/tcp"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

func (f *Factory) buildTCPHandlers(ctx context.Context, route *tcprouter.Router, rtConf *config.Endpoint) error {
	var (
		sHandler tcp.Handler
	)

	for _, router := range rtConf.Routers {
		chain := tcp.NewChain()
		middleware := f.buildTCPMiddleware(ctx, router.Middlewares)

		sHandler = f.buildTCPRouterHandlers(ctx, router)
		if router.Type != config.RuleTypeTCP {
			continue
		}

		then, err := chain.Extend(*middleware).Then(sHandler)
		if err != nil {
			log.Error().Msgf("Error when create tcp router chain, %v", err)
			continue
		}
		if router.TlsEnabled {
			handler := &tcp.TLSHandler{
				Next:   then,
				Config: rtConf.TLSConfig.Config,
			}
			err := route.AddTLSRoute(fmt.Sprintf("HostSNI(`%s`)", router.Host), 0, handler)
			if err != nil {
				log.Error().Msgf("Error When AddRoute, %v", err)
			}
		} else {
			err := route.AddRoute("HostSNI(`*`)", 0, then)
			if err != nil {
				log.Error().Msgf("Error When AddRoute, %v", err)
			}
		}
	}

	return nil
}

func (f *Factory) buildTCPRouterHandlers(ctx context.Context, rtConf config.Router) tcp.Handler {
	handler, err := f.upstreamFactory.BuildTCPUpstreamHandler(ctx, &rtConf.Upstream)
	if err != nil {
		return nil
	}

	return handler
}
