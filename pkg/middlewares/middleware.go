package middlewares

import (
	"api_gateway/pkg/middlewares/auth"
	"context"
	"github.com/pkg/errors"
	"net/http"
)

type MiddlewareCfg interface {
	Validator() error
}

func NewMiddlewareWithType(ctx context.Context, next http.Handler, cfg MiddlewareCfg, mType, name string) (http.Handler, error) {
	switch mType {
	case auth.BasicTypeName:
		return auth.NewBasic(ctx, next, cfg.(*auth.BasicAuth), name)
	}

	return nil, errors.New("Can not Parse middleware type")
}
