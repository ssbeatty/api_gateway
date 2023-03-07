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
	"crypto/tls"
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

		//authConfig := auth.BasicAuth{
		//	Users: []string{
		//		"demo:$apr1$lH3nyBaa$/wCu0V3.1kYdpZPHRbiyv/",
		//	},
		//}

		crt, _ := tls.LoadX509KeyPair("demo.com.pem", "demo.com-key.pem")

		backend.ReloadConfig(dynamic.Message{
			ProviderName: backend.Name(),
			Configuration: map[string]config.Endpoint{
				"tcp-1": {
					Name:       "tcp-1",
					ListenPort: 8080,
					Type:       config.EndpointTypeTCP,
					TLSConfig: config.TLS{
						Config: &tls.Config{
							Certificates: []tls.Certificate{crt},
						},
					},
					Routers: []config.Router{
						{
							Host:       "*",
							TlsEnabled: false,
							Type:       config.RuleTypeTCP,
							Upstream: config.Upstream{
								Type: config.UpstreamTypeServer,
								Paths: []string{
									"127.0.0.1:8088",
									"127.0.0.1:8089",
								},
								LoadBalancerType: loadbalancer.LbRoundRobin,
							},
							//Middlewares: []config.Middleware{
							//	{
							//		Name: "ip allow",
							//		Type: ipallowlist.TypeName,
							//		Config: &ipallowlist.TCPIPAllowList{
							//			SourceRange: []string{
							//				"192.168.50.102",
							//			},
							//		},
							//	},
							//},
						},
						{
							Host:       "api.demo.com",
							TlsEnabled: true,
							Type:       config.RuleTypeTCP,
							Upstream: config.Upstream{
								Type: config.UpstreamTypeServer,
								Paths: []string{
									"api.demo.com:8443",
								},
								LoadBalancerType: loadbalancer.LbRoundRobin,
							},
						},
					},
				},
				//"udp-1": {
				//	Name:       "udp-1",
				//	Type:       config.EndpointTypeUDP,
				//	ListenPort: 30001,
				//	Routers: []config.Router{
				//		{
				//			Type: config.RuleTypeUDP,
				//			Upstream: config.Upstream{
				//				Type: config.UpstreamTypeServer,
				//				Paths: []string{
				//					"0.0.0.0:30000",
				//				},
				//			},
				//		},
				//	},
				//},
			},
		})
	})

	upstreamFactory := upstream.NewFactory(config.DefaultConfig.Gateway, routinesPool)

	routerFactory := routerManager.NewRouterFactory(config.DefaultConfig.Gateway, upstreamFactory)

	// server start
	server := gateway.NewServer(routinesPool, w, routerFactory, &config.DefaultConfig.Gateway)

	server.Start(ctx)
	defer server.Close()

	server.Wait()

}
