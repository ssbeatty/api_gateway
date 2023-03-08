package grpcheaders

import (
	"api_gateway/pkg/logs"
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	TypeName = "HeadersGRPC"
)

type Headers struct {
	CustomRequestHeaders map[string]string `json:"customRequestHeaders,omitempty"`
}

func (b *Headers) Schema() (string, error) {

	return "", nil
}

func New(ctx context.Context, config *Headers, name string) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}
		for old, h := range config.CustomRequestHeaders {
			md.Set(old, h)
		}
		if err := ss.SetHeader(md); err != nil {
			return errors.WithMessage(err, "Grpc SetHeader")
		}
		if err := handler(srv, ss); err != nil {
			log.Error().Str(logs.MiddlewareName, name).Err(err).Send()
			return err
		}
		return nil
	}
}
