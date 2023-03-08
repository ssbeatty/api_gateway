package middlewares

import (
	"api_gateway/pkg/middlewares/addprefix"
	"api_gateway/pkg/middlewares/auth"
	"api_gateway/pkg/middlewares/grpc/grpcheaders"
	"api_gateway/pkg/middlewares/grpc/grpcipallowlist"
	"api_gateway/pkg/middlewares/headers"
	httpIPAllow "api_gateway/pkg/middlewares/ipallowlist"
	"api_gateway/pkg/middlewares/ratelimiter"
	"api_gateway/pkg/middlewares/redirect"
	"api_gateway/pkg/middlewares/replacepath"
	"api_gateway/pkg/middlewares/replacepathregex"
	"api_gateway/pkg/middlewares/retry"
	"api_gateway/pkg/middlewares/stripprefix"
	"api_gateway/pkg/middlewares/stripprefixregex"
	"api_gateway/pkg/middlewares/tcp/inflightconn"
	"api_gateway/pkg/middlewares/tcp/ipallowlist"
	"api_gateway/pkg/tcp"
	"context"
	"google.golang.org/grpc"
	"net/http"
)

type MiddlewareCfg interface {
	Schema() (string, error)
}

type defaultHTTPWrap struct {
}

func (d defaultHTTPWrap) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}

type defaultTCPtWrap struct {
}

func (d defaultTCPtWrap) ServeTCP(conn tcp.WriteCloser) {

}

func NewHTTPMiddlewareWithType(ctx context.Context, next http.Handler, cfg MiddlewareCfg, mType, name string) (http.Handler, error) {
	switch mType {
	case auth.TypeName:
		return auth.NewBasic(ctx, next, cfg.(*auth.BasicAuth), name)
	case addprefix.TypeName:
		return addprefix.New(ctx, next, cfg.(*addprefix.AddPrefix), name)
	case headers.TypeName:
		return headers.NewHeader(next, cfg.(*headers.Headers), name)
	case httpIPAllow.TypeName:
		return httpIPAllow.New(ctx, next, cfg.(*httpIPAllow.IPAllowList), name)
	case ratelimiter.TypeName:
		return ratelimiter.New(ctx, next, cfg.(*ratelimiter.RateLimit), name)
	case redirect.TypeSchemeName:
		return redirect.NewRedirectScheme(ctx, next, cfg.(*redirect.Scheme), name)
	case redirect.TypeRegexName:
		return redirect.NewRedirectRegex(ctx, next, cfg.(*redirect.Regex), name)
	case replacepath.TypeName:
		return replacepath.New(ctx, next, cfg.(*replacepath.ReplacePath), name)
	case replacepathregex.TypeName:
		return replacepathregex.New(ctx, next, cfg.(*replacepathregex.ReplacePathRegex), name)
	case retry.TypeName:
		return retry.New(ctx, next, cfg.(*retry.Retry), retry.Listeners{}, name)
	case stripprefix.TypeName:
		return stripprefix.New(ctx, next, cfg.(*stripprefix.StripPrefix), name)
	case stripprefixregex.TypeName:
		return stripprefixregex.New(ctx, next, cfg.(*stripprefixregex.StripPrefixRegex), name)
	}

	return &defaultHTTPWrap{}, nil
}

func NewTCPMiddlewareWithType(ctx context.Context, next tcp.Handler, cfg MiddlewareCfg, mType, name string) (tcp.Handler, error) {
	switch mType {
	case ipallowlist.TypeName:
		return ipallowlist.New(ctx, next, cfg.(*ipallowlist.TCPIPAllowList), name)
	case inflightconn.TypeName:
		return inflightconn.New(ctx, next, cfg.(*inflightconn.TCPInFlightConn), name)
	}

	return &defaultTCPtWrap{}, nil
}

func NewGRPCMiddlewareWithType(ctx context.Context, cfg MiddlewareCfg, mType, name string) grpc.StreamServerInterceptor {
	switch mType {
	case grpcheaders.TypeName:
		return grpcheaders.New(ctx, cfg.(*grpcheaders.Headers), name)
	case grpcipallowlist.TypeName:
		return grpcipallowlist.New(ctx, cfg.(*grpcipallowlist.IPAllowList), name)
	}
	return nil
}
