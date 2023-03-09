package router

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/pkg/tcp"
	"api_gateway/pkg/types"
	"context"
	"github.com/e421083458/grpc-proxy/proxy"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type GrpcServer struct {
	Server    types.GrpcStop
	Forwarder *GrpcForwarder
}

type GrpcForwarder struct {
	net.Listener
	connChan chan net.Conn
	errChan  chan error
}

func NewGrpcForwarder(ln net.Listener) *GrpcForwarder {
	return &GrpcForwarder{
		Listener: ln,
		connChan: make(chan net.Conn),
		errChan:  make(chan error),
	}
}

// ServeTCP uses the connection to serve it later in "Accept".
func (h *GrpcForwarder) ServeTCP(conn tcp.WriteCloser) {
	h.connChan <- conn
}

// Accept retrieves a served connection in ServeTCP.
func (h *GrpcForwarder) Accept() (net.Conn, error) {
	select {
	case conn := <-h.connChan:
		return conn, nil
	case err := <-h.errChan:
		return nil, err
	}
}

// Close do nothing be.
func (h *GrpcForwarder) Close() error {
	return nil
}

// Error to close listen
func (h *GrpcForwarder) Error(err error) {
	h.errChan <- err
}

func getGrpcRouters(rtConf *config.Endpoint, tls bool) []config.Router {
	var (
		tlsRouters   []config.Router
		notlsRouters []config.Router
	)
	for _, router := range rtConf.Routers {
		if router.TlsEnabled && router.Type == config.RuleTypeGRPC {
			tlsRouters = append(tlsRouters, router)
			continue
		} else if router.Type == config.RuleTypeGRPC {
			notlsRouters = append(notlsRouters, router)
		}
	}

	if tls {
		return tlsRouters
	}

	return notlsRouters
}

func (f *Factory) buildGrpcHandlers(ctx context.Context, rtConf *config.Endpoint, tls bool) *GrpcServer {

	var (
		grpcServers []*GrpcServer
	)

	for _, router := range getGrpcRouters(rtConf, tls) {
		var (
			s *grpc.Server
		)

		grpcHandler, err := f.upstreamFactory.BuildGRPCUpstreamHandler(ctx, &router.Upstream)
		if err != nil {
			log.Error().Msgf("Error when build grpc upstream, %v", err)
			continue
		}

		middlewares := f.buildGrpcMiddleware(ctx, router.Middlewares)

		if tls {
			tlsConfig, err := generateTLSConfig(&router.TLSConfig)
			if err != nil {
				log.Error().Err(err).Msg("Generate tls config")
			}
			c := credentials.NewTLS(tlsConfig)
			s = grpc.NewServer(
				middlewares,
				grpc.Creds(c),
				grpc.CustomCodec(proxy.Codec()),
				grpc.UnknownServiceHandler(grpcHandler),
			)
		} else {
			s = grpc.NewServer(
				middlewares,
				grpc.CustomCodec(proxy.Codec()),
				grpc.UnknownServiceHandler(grpcHandler),
			)
		}

		grpcServers = append(grpcServers, &GrpcServer{
			Server: s,
		})
	}

	if len(grpcServers) > 0 {
		return grpcServers[0]
	}

	return nil
}
