package router

import (
	"api_gateway/internal/gateway/config"
	middleware "api_gateway/pkg/middlewares"
	"api_gateway/pkg/tcp"
	"context"
	"github.com/containous/alice"
	"net/http"
)

func (f *Factory) buildHttpMiddleware(ctx context.Context, middlewares []config.Middleware) *alice.Chain {
	chain := alice.New()
	for _, mid := range middlewares {
		middlewareName := mid.Name
		chain = chain.Append(func(next http.Handler) (http.Handler, error) {
			return middleware.NewHTTPMiddlewareWithType(ctx, next, mid.Config, mid.Type, middlewareName)
		})
	}
	return &chain
}

func (f *Factory) buildTCPMiddleware(ctx context.Context, middlewares []config.Middleware) *tcp.Chain {
	chain := tcp.NewChain()
	for _, mid := range middlewares {
		middlewareName := mid.Name
		chain = chain.Append(func(next tcp.Handler) (tcp.Handler, error) {
			return middleware.NewTCPMiddlewareWithType(ctx, next, mid.Config, mid.Type, middlewareName)
		})
	}
	return &chain
}
