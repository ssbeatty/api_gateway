package main

import (
	"api_gateway/internal/gateway"
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/dynamic"
	routerManager "api_gateway/internal/gateway/manager/router"
	"api_gateway/internal/gateway/manager/upstream"
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
	"api_gateway/internal/gateway/provider"
	"api_gateway/internal/gateway/watcher"

	"api_gateway/pkg/logs"
	"api_gateway/pkg/safe"
	"context"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// init default
	config.SetupConfig()

	logConfiguration := config.DefaultConfig.Log
	logs.SetupLogger(logConfiguration.Level, logConfiguration.Path)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	routinesPool := safe.NewPool(ctx, 10000)

	// backend provider
	backend := provider.NewBackend()
	w := watcher.NewConfigurationWatcher(routinesPool)

	w.AddProvider(backend)

	// todo this is test code
	routinesPool.Go(func() {
		time.Sleep(time.Second * 2)

		backend.ReloadConfig(dynamic.Message{
			ProviderName: backend.Name(),
			Configuration: map[string]config.Endpoint{
				"tcp-1": {
					Name:       "tcp-1",
					ListenPort: 8080,
					Type:       config.EndpointTypeTCP,
					Routers: []config.Routers{
						{
							Host:       "*",
							TlsEnabled: false,
							Rules: []config.Rule{
								{
									Type: config.RuleTypeHTTP,
									Rule: "PathPrefix(`/`)",
									Upstream: config.Upstream{
										Type: config.UpstreamTypeURL,
										Paths: []string{
											"http://127.0.0.1:8088",
											"http://127.0.0.1:8089",
										},
										LoadBalancerType: loadbalancer.LbRoundRobin,
									},
								},
							},
						},
					},
				},
			},
		})
	})

	upstreamFactory := upstream.NewFactory(config.DefaultConfig.Gateway, routinesPool)

	routerFactory := routerManager.NewRouterFactory(config.DefaultConfig.Gateway, upstreamFactory)

	// server start
	server := gateway.NewServer(routinesPool, w, routerFactory)

	server.Start(ctx)
	defer server.Close()

	server.Wait()

}
