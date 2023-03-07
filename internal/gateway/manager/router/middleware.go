package router

import (
	"api_gateway/internal/gateway/config"
	middleware "api_gateway/pkg/middlewares"
	"api_gateway/pkg/tcp"
	"context"
	"github.com/containous/alice"
	"google.golang.org/grpc"
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

func (f *Factory) buildGrpcMiddleware(ctx context.Context, middlewares []config.Middleware) grpc.ServerOption {
	var streamServerInterceptors []grpc.StreamServerInterceptor

	for _, mid := range middlewares {
		middlewareName := mid.Name
		streamServerInterceptor := middleware.NewGRPCMiddlewareWithType(ctx, mid.Config, mid.Type, middlewareName)
		if streamServerInterceptor != nil {
			streamServerInterceptors = append(streamServerInterceptors, streamServerInterceptor)
		}
	}

	return grpc.ChainStreamInterceptor(streamServerInterceptors...)
}
