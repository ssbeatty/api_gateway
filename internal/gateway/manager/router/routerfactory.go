package router

import (
	"api_gateway/internal/gateway/config"
	httpmuxer "api_gateway/internal/gateway/muxer/http"
	tcprouter "api_gateway/internal/gateway/router/tcp"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/middlewares/recovery"
	"context"
	"github.com/containous/alice"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Factory the factory of TCP/UDP routers.
type Factory struct {
	staticConfiguration config.Gateway
}

// NewRouterFactory creates a new RouterFactory.
func NewRouterFactory(staticConfiguration config.Gateway) *Factory {

	return &Factory{
		staticConfiguration: staticConfiguration,
	}
}

// CreateTCPRouters creates new TCPRouter.
func (f *Factory) CreateTCPRouters(ctx context.Context, rtConf *config.Endpoint) *tcprouter.Router {

	handlersNonTLS := buildHttpHandlers(ctx, rtConf, false)
	router, err := tcprouter.NewRouter()
	if err != nil {
		return nil
	}

	router.SetHTTPHandler(handlersNonTLS)
	return router
}

func buildHttpHandlers(ctx context.Context, rtConf *config.Endpoint, tls bool) http.Handler {
	logger := log.With().Str(logs.EndpointName, rtConf.Name).Logger()

	muxer, err := httpmuxer.NewMuxer()
	if err != nil {
		return nil
	}

	// todo build middle
	for _, router := range rtConf.Routers {
		if tls || router.TlsEnabled || router.Host != "*" {
			continue
		}
		for _, rule := range router.Rules {
			handler, buildErr := buildRouterHandler(rule)
			if buildErr != nil {
				logger.Error().Err(err).Send()
				continue
			}

			err = muxer.AddRoute(rule.Rule, rule.Priority, handler)
			if err != nil {
				logger.Error().Err(err).Send()
				continue
			}
		}
	}

	muxer.SortRoutes()
	chain := alice.New()
	chain = chain.Append(func(next http.Handler) (http.Handler, error) {
		return recovery.New(ctx, next)
	})

	newChain, err := chain.Then(muxer)
	if err != nil {
		logger.Error().Err(err).Send()
	}

	return newChain
}

func buildRouterHandler(rule config.Rule) (http.Handler, error) {
	switch rule.Upstream.Type {
	case config.UpstreamTypeURL:
		u, err := url.Parse(rule.Upstream.Paths[0])
		if err != nil {
			return nil, err
		}
		return httputil.NewSingleHostReverseProxy(u), nil
	case config.UpstreamTypeSTATIC:

		return http.FileServer(http.Dir(rule.Upstream.Paths[0])), nil
	case config.UpstreamTypeServer:
	}
	return nil, nil
}
