package router

import (
	"api_gateway/internal/gateway/config"
	middleware "api_gateway/pkg/middlewares"
	"context"
	"github.com/containous/alice"
	"net/http"
)

func (f *Factory) buildHttpMiddleware(ctx context.Context, middlewares []config.Middleware) *alice.Chain {
	chain := alice.New()
	for _, mid := range middlewares {
		middlewareName := mid.Name
		chain = chain.Append(func(next http.Handler) (http.Handler, error) {
			return middleware.NewMiddlewareWithType(ctx, next, mid.Config, mid.Type, middlewareName)
		})
	}
	return &chain
}
