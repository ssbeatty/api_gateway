package udp

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/pkg/safe"
	"api_gateway/pkg/udp"
	"context"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

// EndPoint is an entry point where we listen for UDP packets.
type EndPoint struct {
	listener               *udp.Listener
	switcher               *udp.HandlerSwitcher
	transportConfiguration *config.Gateway
	pool                   *safe.Pool
}

// NewUDPEntryPoint returns a UDP entry point.
func NewUDPEntryPoint(cfg *config.Endpoint, transportConfiguration *config.Gateway, pool *safe.Pool) (*EndPoint, error) {
	addr, err := net.ResolveUDPAddr("udp", cfg.GetAddress())
	if err != nil {
		return nil, err
	}

	listener, err := udp.Listen("udp", addr, transportConfiguration.ListenUDPTimeOut)
	if err != nil {
		return nil, err
	}

	return &EndPoint{
		listener:               listener,
		switcher:               &udp.HandlerSwitcher{},
		transportConfiguration: transportConfiguration,
		pool:                   pool,
	}, nil
}

// Start commences the listening for ep.
func (ep *EndPoint) Start(ctx context.Context) {
	log.Ctx(ctx).Debug().Msg("Start UDP Server")
	for {
		conn, err := ep.listener.Accept()
		if err != nil {
			// Only errClosedListener can happen that's why we return
			return
		}

		ep.pool.Go(func() {
			ep.switcher.ServeUDP(conn)
		})
	}
}

// Shutdown closes ep's listener. It eventually closes all "sessions" and
// releases associated resources, but only after it has waited for a graceTimeout,
// if any was configured.
func (ep *EndPoint) Shutdown(ctx context.Context) {
	logger := log.Ctx(ctx)

	reqAcceptGraceTimeOut := ep.transportConfiguration.GraceTimeOut
	if reqAcceptGraceTimeOut > 0 {
		logger.Info().Msgf("Waiting %s for incoming requests to cease", reqAcceptGraceTimeOut)
		time.Sleep(reqAcceptGraceTimeOut)
	}

	graceTimeOut := ep.transportConfiguration.GraceTimeOut
	if err := ep.listener.Shutdown(graceTimeOut); err != nil {
		logger.Error().Err(err).Send()
	}
}

// Switch replaces ep's handler with the one given as argument.
func (ep *EndPoint) Switch(handler udp.Handler) {
	ep.switcher.Switch(handler)
}
