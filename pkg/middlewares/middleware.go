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
	"encoding/json"
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

func unmarshalAny[T any](cfg map[string]interface{}) (*T, error) {
	marshal, err := json.Marshal(&cfg)
	if err != nil {
		return nil, err
	}

	out := new(T)
	if err := json.Unmarshal(marshal, out); err != nil {
		return nil, err
	}
	return out, nil
}

func NewHTTPMiddlewareWithType(ctx context.Context, next http.Handler, cfg map[string]interface{}, mType, name string) (http.Handler, error) {

	switch mType {
	case auth.TypeName:
		if c, e := unmarshalAny[auth.BasicAuth](cfg); e != nil {
			return nil, e
		} else {
			return auth.NewBasic(ctx, next, c, name)
		}
	case addprefix.TypeName:
		if c, e := unmarshalAny[addprefix.AddPrefix](cfg); e != nil {
			return nil, e
		} else {
			return addprefix.New(ctx, next, c, name)
		}
	case headers.TypeName:
		if c, e := unmarshalAny[headers.Headers](cfg); e != nil {
			return nil, e
		} else {
			return headers.NewHeader(next, c, name)
		}
	case httpIPAllow.TypeName:
		if c, e := unmarshalAny[httpIPAllow.IPAllowList](cfg); e != nil {
			return nil, e
		} else {
			return httpIPAllow.New(ctx, next, c, name)
		}
	case ratelimiter.TypeName:
		if c, e := unmarshalAny[ratelimiter.RateLimit](cfg); e != nil {
			return nil, e
		} else {
			return ratelimiter.New(ctx, next, c, name)
		}
	case redirect.TypeSchemeName:
		if c, e := unmarshalAny[redirect.Scheme](cfg); e != nil {
			return nil, e
		} else {
			return redirect.NewRedirectScheme(ctx, next, c, name)
		}
	case redirect.TypeRegexName:
		if c, e := unmarshalAny[redirect.Regex](cfg); e != nil {
			return nil, e
		} else {
			return redirect.NewRedirectRegex(ctx, next, c, name)
		}
	case replacepath.TypeName:
		if c, e := unmarshalAny[replacepath.ReplacePath](cfg); e != nil {
			return nil, e
		} else {
			return replacepath.New(ctx, next, c, name)
		}
	case replacepathregex.TypeName:
		if c, e := unmarshalAny[replacepathregex.ReplacePathRegex](cfg); e != nil {
			return nil, e
		} else {
			return replacepathregex.New(ctx, next, c, name)
		}
	case retry.TypeName:
		if c, e := unmarshalAny[retry.Retry](cfg); e != nil {
			return nil, e
		} else {
			return retry.New(ctx, next, c, retry.Listeners{}, name)
		}
	case stripprefix.TypeName:
		if c, e := unmarshalAny[stripprefix.StripPrefix](cfg); e != nil {
			return nil, e
		} else {
			return stripprefix.New(ctx, next, c, name)
		}
	case stripprefixregex.TypeName:
		if c, e := unmarshalAny[stripprefixregex.StripPrefixRegex](cfg); e != nil {
			return nil, e
		} else {
			return stripprefixregex.New(ctx, next, c, name)
		}
	}

	return &defaultHTTPWrap{}, nil
}

func NewTCPMiddlewareWithType(ctx context.Context, next tcp.Handler, cfg map[string]interface{}, mType, name string) (tcp.Handler, error) {
	switch mType {
	case ipallowlist.TypeName:
		if c, e := unmarshalAny[ipallowlist.TCPIPAllowList](cfg); e != nil {
			return nil, e
		} else {
			return ipallowlist.New(ctx, next, c, name)
		}
	case inflightconn.TypeName:
		if c, e := unmarshalAny[inflightconn.TCPInFlightConn](cfg); e != nil {
			return nil, e
		} else {
			return inflightconn.New(ctx, next, c, name)
		}
	}

	return &defaultTCPtWrap{}, nil
}

func NewGRPCMiddlewareWithType(ctx context.Context, cfg map[string]interface{}, mType, name string) grpc.StreamServerInterceptor {
	switch mType {
	case grpcheaders.TypeName:
		if c, e := unmarshalAny[grpcheaders.Headers](cfg); e != nil {
			return nil
		} else {
			return grpcheaders.New(ctx, c, name)
		}
	case grpcipallowlist.TypeName:
		if c, e := unmarshalAny[grpcipallowlist.IPAllowList](cfg); e != nil {
			return nil
		} else {
			return grpcipallowlist.New(ctx, c, name)
		}
	}
	return nil
}
