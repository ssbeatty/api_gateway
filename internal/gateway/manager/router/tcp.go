package router

import (
	"api_gateway/internal/gateway/config"
	tcprouter "api_gateway/internal/gateway/router/tcp"
	"api_gateway/pkg/tcp"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (f *Factory) buildTCPHandlers(ctx context.Context, route *tcprouter.Router, rtConf *config.Endpoint) (tcp.Handler, error) {
	var (
		sHandler tcp.Handler
	)

	for _, router := range rtConf.Routers {
		if len(router.Rules) != 1 {
			return nil, errors.New("TCP route has only one rule")
		}
		rule := &router.Rules[0]
		sHandler = f.buildTCPRouterHandlers(ctx, rule)
		if rule.Type != config.RuleTypeTCP {
			continue
		}

		if router.TlsEnabled {
			err := route.AddTLSRoute(fmt.Sprintf("HostSNI(`%s`)", router.Host), 0, sHandler)
			if err != nil {
				log.Error().Msgf("Error When AddRoute, %v", err)
			}
		} else {
			err := route.AddRoute("HostSNI(`*`)", 0, sHandler)
			if err != nil {
				log.Error().Msgf("Error When AddRoute, %v", err)
			}
		}
	}

	return tcp.NewChain().Then(sHandler)
}

func (f *Factory) buildTCPRouterHandlers(ctx context.Context, rtConf *config.Rule) tcp.Handler {
	handler, err := f.upstreamFactory.BuildTCPUpstreamHandler(ctx, &rtConf.Upstream)
	if err != nil {
		return nil
	}

	return handler
}
