package upstream

import (
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
	"api_gateway/pkg/udp"
	"context"
	"github.com/rs/zerolog/log"
)

func (f *Factory) NewUDPLoadBalanceReverseProxy(ctx context.Context, lb loadbalancer.LoadBalance) *UDPReverseProxy {
	return func() *UDPReverseProxy {

		return &UDPReverseProxy{
			ctx: ctx,
			lb:  lb,
		}
	}()
}

type UDPReverseProxy struct {
	ctx context.Context
	lb  loadbalancer.LoadBalance
}

// ServeUDP forwards the connection to the right service.
func (b *UDPReverseProxy) ServeUDP(conn *udp.Conn) {
	nextAddr, err := b.lb.Get("")
	if err != nil {
		log.Error().Err(err).Send()
	}

	handler, err := udp.NewProxy(nextAddr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create server")
		return
	}

	handler.ServeUDP(conn)

}
