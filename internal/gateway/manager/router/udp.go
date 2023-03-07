package router

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/pkg/udp"
	"github.com/rs/zerolog/log"
)

func (f *Factory) buildUDPHandlers(rtConf *config.Endpoint) []udp.Handler {
	var handlers []udp.Handler

	for _, router := range rtConf.Routers {
		loadBalancer := udp.NewWRRLoadBalancer()

		if router.Type != config.RuleTypeUDP {
			continue
		}
		for _, address := range router.Upstream.Paths {
			handler, err := udp.NewProxy(address)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create server")
				continue
			}
			loadBalancer.AddServer(handler)
		}

		handlers = append(handlers, loadBalancer)

	}

	return handlers
}
