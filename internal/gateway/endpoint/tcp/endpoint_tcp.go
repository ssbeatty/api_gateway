package tcp

import (
	"api_gateway/internal/gateway/config"
	routerManager "api_gateway/internal/gateway/manager/router"
	"api_gateway/internal/gateway/muxer/requestdecorator"
	"api_gateway/internal/gateway/router"
	tcprouter "api_gateway/internal/gateway/router/tcp"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/middlewares/contenttype"
	"api_gateway/pkg/middlewares/forwardedheaders"
	"api_gateway/pkg/safe"
	"api_gateway/pkg/tcp"
	"context"
	"errors"
	"fmt"
	"github.com/containous/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	stdlog "log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"syscall"
	"time"
)

type httpForwarder struct {
	net.Listener
	connChan chan net.Conn
	errChan  chan error
}

func newHTTPForwarder(ln net.Listener) *httpForwarder {
	return &httpForwarder{
		Listener: ln,
		connChan: make(chan net.Conn),
		errChan:  make(chan error),
	}
}

// ServeTCP uses the connection to serve it later in "Accept".
func (h *httpForwarder) ServeTCP(conn tcp.WriteCloser) {
	h.connChan <- conn
}

// Accept retrieves a served connection in ServeTCP.
func (h *httpForwarder) Accept() (net.Conn, error) {
	select {
	case conn := <-h.connChan:
		return conn, nil
	case err := <-h.errChan:
		return nil, err
	}
}

type EndPoint struct {
	listener      net.Listener
	switcher      *tcp.HandlerSwitcher
	tracker       *connectionTracker
	httpServer    *httpServer
	httpsServer   *httpServer
	grpcServer    *routerManager.GrpcServer
	pool          *safe.Pool
	configuration *config.Endpoint
}

// NewTCPEndPoint creates a new TCPEndPoint.
func NewTCPEndPoint(ctx context.Context, configuration *config.Endpoint, pool *safe.Pool) (*EndPoint, error) {
	tracker := newConnectionTracker()

	listener, err := buildListener(configuration)
	if err != nil {
		return nil, fmt.Errorf("error preparing server: %w", err)
	}

	rt := &tcprouter.Router{}
	reqDecorator := requestdecorator.New()

	httpServer, err := createHTTPServer(ctx, listener, true, reqDecorator)
	if err != nil {
		return nil, fmt.Errorf("error preparing http server: %w", err)
	}

	rt.SetHTTPForwarder(httpServer.Forwarder)

	httpsServer, err := createHTTPServer(ctx, listener, false, reqDecorator)
	if err != nil {
		return nil, fmt.Errorf("error preparing https server: %w", err)
	}

	rt.SetHTTPSForwarder(httpsServer.Forwarder)

	tcpSwitcher := &tcp.HandlerSwitcher{}
	tcpSwitcher.Switch(rt)

	return &EndPoint{
		listener:      listener,
		switcher:      tcpSwitcher,
		tracker:       tracker,
		httpServer:    httpServer,
		httpsServer:   httpsServer,
		pool:          pool,
		configuration: configuration,
	}, nil
}

// Start starts the TCP server.
func (e *EndPoint) Start(ctx context.Context) {
	logger := log.Ctx(ctx)
	logger.Debug().Msg("Starting TCP Server")

	for {
		conn, err := e.listener.Accept()
		if err != nil {
			logger.Error().Err(err).Send()

			var opErr *net.OpError
			if errors.As(err, &opErr) && opErr.Temporary() {
				continue
			}

			var urlErr *url.Error
			if errors.As(err, &urlErr) && urlErr.Temporary() {
				continue
			}

			e.httpServer.Forwarder.errChan <- err
			e.httpsServer.Forwarder.errChan <- err

			return
		}

		writeCloser, err := writeCloser(conn)
		if err != nil {
			panic(err)
		}

		e.pool.Go(func() {
			e.switcher.ServeTCP(newTrackedConnection(writeCloser, e.tracker))
		})
	}
}

// Shutdown stops the TCP connections.
func (e *EndPoint) Shutdown(ctx context.Context) {
	logger := log.Ctx(ctx)

	var (
		cancel context.CancelFunc
	)

	graceTimeOut := config.DefaultConfig.Gateway.GraceTimeOut
	if config.DefaultConfig.Gateway.GraceTimeOut > 0 {
		_, cancel = context.WithTimeout(ctx, graceTimeOut)
		logger.Debug().Msgf("Waiting %s seconds before killing connections", graceTimeOut)
	}

	var wg sync.WaitGroup

	shutdownServer := func(server stoppable) {
		defer wg.Done()
		err := server.Shutdown(ctx)
		if err == nil {
			return
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			logger.Debug().Err(err).Msg("Server failed to shutdown within deadline")
			if err = server.Close(); err != nil {
				logger.Error().Err(err).Send()
			}
			return
		}

		logger.Error().Err(err).Send()

		// We expect Close to fail again because Shutdown most likely failed when trying to close a listener.
		// We still call it however, to make sure that all connections get closed as well.
		server.Close()
	}

	if e.httpServer.Server != nil {
		wg.Add(1)
		go shutdownServer(e.httpServer.Server)
	}

	if e.httpsServer.Server != nil {
		wg.Add(1)
		go shutdownServer(e.httpsServer.Server)

	}

	if e.tracker != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := e.tracker.Shutdown(ctx)
			if err == nil {
				return
			}
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				logger.Debug().Err(err).Msg("Server failed to shutdown before deadline")
			}
			e.tracker.Close()
		}()
	}

	wg.Wait()
	cancel()
}

// SwitchRouter switches the TCP router handler.
func (e *EndPoint) SwitchRouter(rt *tcprouter.Router, gs *routerManager.GrpcServer) {
	rt.SetHTTPForwarder(e.httpServer.Forwarder)

	httpHandler := rt.GetHTTPHandler()
	if httpHandler == nil {
		httpHandler = router.BuildDefaultHTTPRouter()
	}

	e.httpServer.Switcher.UpdateHandler(httpHandler)

	rt.SetHTTPSForwarder(e.httpsServer.Forwarder)

	httpsHandler := rt.GetHTTPSHandler()
	if httpsHandler == nil {
		httpsHandler = router.BuildDefaultHTTPRouter()
	}

	e.httpsServer.Switcher.UpdateHandler(httpsHandler)

	e.switcher.Switch(rt)
	e.grpcServer = gs

	if gs != nil {

		gs.Forwarder = routerManager.NewGrpcForwarder(e.listener)
		rt.SetGRPCForwarder(gs.Forwarder)
		err := e.grpcServer.Server.Serve(gs.Forwarder)
		if err != nil {
			log.Error().Err(err).Send()
		}
	}

}

// writeCloser returns the given connection, augmented with the WriteCloser
// implementation, if any was found within the underlying conn.
func writeCloser(conn net.Conn) (tcp.WriteCloser, error) {
	switch typedConn := conn.(type) {
	case *net.TCPConn:
		return typedConn, nil
	default:
		return nil, fmt.Errorf("unknown connection type %T", typedConn)
	}
}

type stoppable interface {
	Shutdown(context.Context) error
	Close() error
}

type stoppableServer interface {
	stoppable
	Serve(listener net.Listener) error
}

type httpServer struct {
	Server    stoppableServer
	Forwarder *httpForwarder
	Switcher  *HTTPHandlerSwitcher
}

func createHTTPServer(ctx context.Context, ln net.Listener, withH2c bool, reqDecorator *requestdecorator.RequestDecorator) (*httpServer, error) {
	httpSwitcher := NewHandlerSwitcher(router.BuildDefaultHTTPRouter())

	next, err := alice.New(requestdecorator.WrapHandler(reqDecorator)).Then(httpSwitcher)
	if err != nil {
		return nil, err
	}

	var handler http.Handler
	// todo
	handler, err = forwardedheaders.NewXForwarded(
		true,
		nil,
		next)
	if err != nil {
		return nil, err
	}

	handler = http.AllowQuerySemicolons(handler)

	handler = contenttype.DisableAutoDetection(handler)

	if withH2c {
		handler = h2c.NewHandler(handler, &http2.Server{
			MaxConcurrentStreams: uint32(250),
		})
	}

	// todo configs this
	serverHTTP := &http.Server{
		Handler:      handler,
		ErrorLog:     stdlog.New(logs.NoLevel(log.Logger, zerolog.DebugLevel), "", 0),
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}

	listener := newHTTPForwarder(ln)
	go func() {
		err := serverHTTP.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Ctx(ctx).Error().Err(err).Msg("Error while starting server")
		}
	}()
	return &httpServer{
		Server:    serverHTTP,
		Forwarder: listener,
		Switcher:  httpSwitcher,
	}, nil
}

func buildListener(entryPoint *config.Endpoint) (net.Listener, error) {
	listener, err := net.Listen("tcp", entryPoint.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("error opening listener: %w", err)
	}

	listener = tcpKeepAliveListener{listener.(*net.TCPListener)}

	return listener, nil
}

func newConnectionTracker() *connectionTracker {
	return &connectionTracker{
		conns: make(map[net.Conn]struct{}),
	}
}

type connectionTracker struct {
	conns map[net.Conn]struct{}
	lock  sync.RWMutex
}

// AddConnection add a connection in the tracked connections list.
func (c *connectionTracker) AddConnection(conn net.Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.conns[conn] = struct{}{}
}

// RemoveConnection remove a connection from the tracked connections list.
func (c *connectionTracker) RemoveConnection(conn net.Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.conns, conn)
}

func (c *connectionTracker) isEmpty() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.conns) == 0
}

// Shutdown wait for the connection closing.
func (c *connectionTracker) Shutdown(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		if c.isEmpty() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// Close all the connections in the tracked connections list.
func (c *connectionTracker) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for conn := range c.conns {
		if err := conn.Close(); err != nil {
			log.Error().Err(err).Msg("Error while closing connection")
		}
		delete(c.conns, conn)
	}
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}

	if err := tc.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if err := tc.SetKeepAlivePeriod(3 * time.Minute); err != nil {
		// Some systems, such as OpenBSD, have no user-settable per-socket TCP keepalive options.
		if !errors.Is(err, syscall.ENOPROTOOPT) {
			return nil, err
		}
	}

	return tc, nil
}

func newTrackedConnection(conn tcp.WriteCloser, tracker *connectionTracker) *trackedConnection {
	tracker.AddConnection(conn)
	return &trackedConnection{
		WriteCloser: conn,
		tracker:     tracker,
	}
}

type trackedConnection struct {
	tracker *connectionTracker
	tcp.WriteCloser
}

func (t *trackedConnection) Close() error {
	t.tracker.RemoveConnection(t.WriteCloser)
	return t.WriteCloser.Close()
}
