package router

import (
	"api_gateway/internal/gateway/config"
	httpmuxer "api_gateway/internal/gateway/muxer/http"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/middlewares/recovery"
	"context"
	"github.com/containous/alice"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

func getRouters(rtConf *config.Endpoint, tls bool) []config.Router {
	var (
		tlsRouters   []config.Router
		notlsRouters []config.Router
	)
	for _, router := range rtConf.Routers {
		if router.TlsEnabled {
			tlsRouters = append(tlsRouters, router)
			continue
		}
		notlsRouters = append(notlsRouters, router)
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

	for _, router := range getRouters(rtConf, tls) {
		if router.Type != config.RuleTypeHTTP {
			continue
		}
		// create every rule middleware for everyone http handler
		// example /handler1 has auth middleware but /handler2 not
		middleware := f.buildHttpMiddleware(ctx, router.Middlewares)

		handler, buildErr := f.buildHttpRouterHandler(router)
		if buildErr != nil {
			logger.Debug().Msgf("Build http router error, %v", buildErr)
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

	newChain, err := chain.Then(muxer)
	if err != nil {
		logger.Error().Err(err).Send()
	}

	return newChain
}

func (f *Factory) buildHttpRouterHandler(rule config.Router) (http.Handler, error) {
	if len(rule.Upstream.Paths) == 0 {
		return nil, errors.New("Empty Services!")
	}
	switch rule.Upstream.Type {
	case config.UpstreamTypeURL:
		return f.upstreamFactory.BuildHttpUpstreamHandler(&rule.Upstream)
	case config.UpstreamTypeSTATIC:
		return http.FileServer(http.Dir(rule.Upstream.Paths[0])), nil
	case config.UpstreamTypeServer:
		return nil, errors.New("Http Handlers Can't Be Server")
	}

	return nil, errors.New("Upstream Has not Type")
}
