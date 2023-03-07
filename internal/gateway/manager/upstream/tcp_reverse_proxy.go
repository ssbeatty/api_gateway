package upstream

import (
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
	"api_gateway/pkg/tcp"
	"context"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"time"
)

func (f *Factory) NewTcpLoadBalanceReverseProxy(ctx context.Context, lb loadbalancer.LoadBalance) *TcpReverseProxy {
	return func() *TcpReverseProxy {

		return &TcpReverseProxy{
			ctx:             ctx,
			KeepAlivePeriod: time.Second,
			DialTimeout:     time.Second,
			lb:              lb,
		}
	}()
}

type TcpReverseProxy struct {
	ctx                  context.Context
	Addr                 string
	KeepAlivePeriod      time.Duration
	DialTimeout          time.Duration
	DialContext          func(ctx context.Context, network, address string) (net.Conn, error)
	OnDialError          func(src net.Conn, dstDialErr error)
	ProxyProtocolVersion int
	lb                   loadbalancer.LoadBalance
}

func (dp *TcpReverseProxy) dialTimeout() time.Duration {
	if dp.DialTimeout > 0 {
		return dp.DialTimeout
	}
	return 10 * time.Second
}

func (dp *TcpReverseProxy) dialContext() func(ctx context.Context, network, address string) (net.Conn, error) {
	if dp.DialContext != nil {
		return dp.DialContext
	}
	return (&net.Dialer{
		Timeout:   dp.DialTimeout,
		KeepAlive: dp.KeepAlivePeriod,
	}).DialContext
}

func (dp *TcpReverseProxy) keepAlivePeriod() time.Duration {
	if dp.KeepAlivePeriod != 0 {
		return dp.KeepAlivePeriod
	}
	return time.Minute
}

func (dp *TcpReverseProxy) ServeTCP(src tcp.WriteCloser) {

	var (
		cancel context.CancelFunc
		ctx    context.Context
	)
	if dp.DialTimeout >= 0 {
		ctx, cancel = context.WithTimeout(dp.ctx, dp.dialTimeout())
	}

	nextAddr, err := dp.lb.Get("")
	if err != nil {
		log.Error().Err(err).Send()
	}

	dst, err := dp.dialContext()(ctx, "tcp", nextAddr)
	if cancel != nil {
		cancel()
	}
	if err != nil {
		dp.onDialError()(src, err)
		return
	}

	defer dst.Close()

	if ka := dp.keepAlivePeriod(); ka > 0 {
		if c, ok := dst.(*net.TCPConn); ok {
			_ = c.SetKeepAlive(true)
			_ = c.SetKeepAlivePeriod(ka)
		}
	}
	errChan := make(chan error, 1)
	go dp.proxyCopy(errChan, src, dst)
	go dp.proxyCopy(errChan, dst, src)
	<-errChan
}

func (dp *TcpReverseProxy) onDialError() func(src net.Conn, dstDialErr error) {
	if dp.OnDialError != nil {
		return dp.OnDialError
	}
	return func(src net.Conn, dstDialErr error) {
		log.Error().Msgf("proxy copy incoming conn %v, error dialing %q: %v", src.RemoteAddr().String(), dp.Addr, dstDialErr)
		src.Close()
	}
}

func (dp *TcpReverseProxy) proxyCopy(errChan chan<- error, dst, src net.Conn) {
	_, err := io.Copy(dst, src)
	errChan <- err
}
