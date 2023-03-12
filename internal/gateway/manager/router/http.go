package router

import (
	"api_gateway/internal/gateway/config"
	httpmuxer "api_gateway/internal/gateway/muxer/http"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/middlewares/accesslog"
	"api_gateway/pkg/middlewares/recovery"
	"api_gateway/pkg/tcp"
	"context"
	"github.com/containous/alice"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
)

type HttpForwarder struct {
	net.Listener
	connChan chan net.Conn
	errChan  chan error
}

func NewHTTPForwarder(ln net.Listener) *HttpForwarder {
	return &HttpForwarder{
		Listener: ln,
		connChan: make(chan net.Conn),
		errChan:  make(chan error),
	}
}

// ServeTCP uses the connection to serve it later in "Accept".
func (h *HttpForwarder) ServeTCP(conn tcp.WriteCloser) {
	h.connChan <- conn
}

// Accept retrieves a served connection in ServeTCP.
func (h *HttpForwarder) Accept() (net.Conn, error) {
	select {
	case conn := <-h.connChan:
		return conn, nil
	case err := <-h.errChan:
		return nil, err
	}
}

// Error to close listen
func (h *HttpForwarder) Error(err error) {
	h.errChan <- err
}

func getHTTPRouters(rtConf *config.Endpoint, tls bool) []config.Router {
	var (
		tlsRouters   []config.Router
		notlsRouters []config.Router
	)
	for _, router := range rtConf.Routers {
		if router.TlsEnabled && router.Type == config.RuleTypeHTTP {
			tlsRouters = append(tlsRouters, router)
			continue
		} else if router.Type == config.RuleTypeHTTP {
			notlsRouters = append(notlsRouters, router)
		}
	}

	if tls {
		return tlsRouters
	}

	return notlsRouters
}

func (f *Factory) buildHttpHandlers(ctx context.Context, rtConf *config.Endpoint, tls bool) http.Handler {
	logger := log.With().Str(logs.EndpointName, rtConf.Name).Logger()

	muxer, err := httpmuxer.NewMuxer()
	if err != nil {
		return nil
	}

	for _, router := range getHTTPRouters(rtConf, tls) {
		// create every rule middleware for everyone http handler
		// example /handler1 has auth middleware but /handler2 not
		middleware := f.buildHttpMiddleware(ctx, router.Middlewares)

		handler, buildErr := f.buildHttpRouter(ctx, router)
		if buildErr != nil {
			logger.Error().Msgf("Build http router error, %v", buildErr)
			continue
		}
		then, chainErr := middleware.Then(handler)
		if chainErr != nil {
			logger.Error().Err(chainErr).Send()
			continue
		}

		err = muxer.AddRoute(router.Rule, router.Priority, then)
		if err != nil {
			logger.Error().Err(err).Send()
			continue
		}
	}

	muxer.SortRoutes()
	chain := alice.New()
	chain = chain.Append(func(next http.Handler) (http.Handler, error) {
		return recovery.New(ctx, next)
	})

	accessConfig := config.DefaultConfig.Log.AccessLog
	if accessConfig.Enable {
		chain = chain.Append(func(next http.Handler) (http.Handler, error) {
			return accesslog.NewHandler(accessConfig, next)
		})
	}

	newChain, err := chain.Then(muxer)
	if err != nil {
		logger.Error().Err(err).Send()
	}

	return newChain
}

func (f *Factory) buildHttpRouter(ctx context.Context, rule config.Router) (http.Handler, error) {
	if len(rule.Upstream.Paths) == 0 {
		return nil, errors.New("Empty Services!")
	}
	switch rule.Upstream.Type {
	case config.UpstreamTypeURL:
		return f.upstreamFactory.BuildHttpUpstreamHandler(ctx, &rule.Upstream)
	case config.UpstreamTypeSTATIC:
		// todo only one path
		return http.FileServer(http.Dir(rule.Upstream.Paths[0])), nil
	case config.UpstreamTypeServer:
		return nil, errors.New("Http Handlers Can't Be Server")
	}

	return nil, errors.New("Upstream Has not Type")
}
