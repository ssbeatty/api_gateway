package router

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/pkg/udp"
	"context"
)

func (f *Factory) buildUDPHandlers(ctx context.Context, rtConf *config.Endpoint) []udp.Handler {
	var handlers []udp.Handler

	for _, router := range rtConf.Routers {
		loadBalancer := f.upstreamFactory.BuildUDPUpstreamHandler(ctx, &router.Upstream)
		handlers = append(handlers, loadBalancer)
	}

	return handlers
}
