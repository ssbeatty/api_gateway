package upstream

import (
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
	"context"
	"github.com/mwitkow/grpc-proxy/proxy"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func (f *Factory) NewGrpcLoadBalanceHandler(lb loadbalancer.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Error().Msgf("Load Balance Poll is empty")
		}
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			if strings.HasPrefix(fullMethodName, "/com.example.internal.") {
				return ctx, nil, status.Errorf(codes.Unimplemented,
					"Unknown method")
			}
			c, err := grpc.DialContext(
				ctx,
				nextAddr,
				grpc.WithCodec(proxy.Codec()),
				grpc.WithInsecure(),
			)
			md, _ := metadata.FromIncomingContext(ctx)
			outCtx, _ := context.WithCancel(ctx)
			outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
			return outCtx, c, err
		}
		return proxy.TransparentHandler(director)
	}()
}
