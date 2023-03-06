package gateway

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/internal/gateway/endpoint/tcp"
	httpmuxer "api_gateway/internal/gateway/muxer/http"
	tcprouter "api_gateway/internal/gateway/router/tcp"
	"api_gateway/internal/gateway/watcher"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/middlewares/recovery"
	"api_gateway/pkg/safe"
	"context"
	"github.com/containous/alice"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type Server struct {
	watcher      *watcher.ConfigurationWatcher
	routinesPool *safe.Pool
	endpoint     *tcp.EndPoint
	stopChan     chan struct{}
}

func NewServer(routinesPool *safe.Pool, watcher *watcher.ConfigurationWatcher) *Server {

	return &Server{
		watcher:      watcher,
		routinesPool: routinesPool,
		stopChan:     make(chan struct{}, 1),
	}
}

// Start starts the server and Stop/Close it when context is Done.
func (s *Server) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		logger := log.Ctx(ctx)
		logger.Info().Msg("I have to go...")
		logger.Info().Msg("Stopping server gracefully")
		s.Stop()
	}()

	s.setupConfigWatcher()
	s.watcher.Start()
}

// Wait blocks until the server shutdown.
func (s *Server) Wait() {
	<-s.stopChan
}

// Stop stops the server.
func (s *Server) Stop() {
	defer log.Info().Msg("Server stopped")

	s.stopChan <- struct{}{}
}

// Close destroys the server.
func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	go func(ctx context.Context) {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			panic("Timeout while stopping, killing instance ✝")
		}
	}(ctx)

	s.routinesPool.Stop()

	cancel()
}

func (s *Server) setupConfigWatcher() {
	s.watcher.AddListener(func(configuration dynamic.Configuration) {
		switch configuration.Action {
		case watcher.ActionCreate:

			// todo test
			s.endpoint, _ = tcp.NewTCPEndPoint(context.Background(), &configuration.EndPoint, s.routinesPool)
			s.endpoint.Start(context.Background())
			switchRouter(s.endpoint)
		}
	})
}

func switchRouter(serverEntryPointsTCP *tcp.EndPoint) func(conf dynamic.Configuration) {
	return func(conf dynamic.Configuration) {
		routers := CreateRouters(&conf)

		serverEntryPointsTCP.SwitchRouter(routers)
	}
}

// CreateRouters creates new TCPRouters and UDPRouters.
func CreateRouters(rtConf *dynamic.Configuration) *tcprouter.Router {
	ctx := context.Background()

	handlersNonTLS := buildHttpHandlers(ctx, rtConf)
	router, err := tcprouter.NewRouter()
	if err != nil {
		return nil
	}

	router.SetHTTPHandler(handlersNonTLS)
	return router
}

func buildHttpHandlers(ctx context.Context, rtConf *dynamic.Configuration) http.Handler {
	logger := log.With().Str(logs.EndpointName, rtConf.EndPoint.Name).Logger()

	muxer, err := httpmuxer.NewMuxer()
	if err != nil {
		return nil
	}

	for _, rule := range rtConf.EndPoint.Rules {
		handler, err := buildRouterHandler(&rule)
		if err != nil {
			logger.Error().Err(err).Send()
			continue
		}

		err = muxer.AddRoute(rule.Router.Rule, rule.Router.Priority, handler)
		if err != nil {
			logger.Error().Err(err).Send()
			continue
		}
	}

	muxer.SortRoutes()
	chain := alice.New()
	chain = chain.Append(func(next http.Handler) (http.Handler, error) {
		return recovery.New(ctx, next)
	})

	newChain, err := chain.Then(muxer)
	if err != nil {
		logger.Error().Err(err).Send()
	}

	return newChain
}

func buildRouterHandler(rule *config.Rule) (http.Handler, error) {

	return nil, nil
}
