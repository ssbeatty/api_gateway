package middlewares

import (
	"api_gateway/pkg/middlewares/auth"
	"api_gateway/pkg/middlewares/ipallowlist"
	"api_gateway/pkg/tcp"
	"context"
	"google.golang.org/grpc"
	"net/http"
)

type MiddlewareCfg interface {
	Validator() error
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
	case auth.BasicTypeName:
		return auth.NewBasic(ctx, next, cfg.(*auth.BasicAuth), name)
	}

	return &defaultHTTPWrap{}, nil
}

func NewTCPMiddlewareWithType(ctx context.Context, next tcp.Handler, cfg MiddlewareCfg, mType, name string) (tcp.Handler, error) {
	switch mType {
	case ipallowlist.TypeName:
		return ipallowlist.New(ctx, next, cfg.(*ipallowlist.TCPIPAllowList), name)
	}

	return &defaultTCPtWrap{}, nil
}

func NewGRPCMiddlewareWithType(ctx context.Context, cfg MiddlewareCfg, mType, name string) grpc.StreamServerInterceptor {
	switch mType {
	case ipallowlist.TypeName:

	}

	return nil
}
