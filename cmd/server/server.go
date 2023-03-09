package main

import (
	"api_gateway/internal/gateway"
	"api_gateway/internal/gateway/config"
	routerManager "api_gateway/internal/gateway/manager/router"
	"api_gateway/internal/gateway/manager/upstream"
	backendProvider "api_gateway/internal/gateway/provider/backend"
	fileProvider "api_gateway/internal/gateway/provider/file"
	"api_gateway/internal/gateway/watcher"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/safe"
	"context"
	"os/signal"
	"syscall"
)

func main() {
	// init default
	config.SetupConfig()

	logConfiguration := config.DefaultConfig.Log
	logs.SetupLogger(logConfiguration.Level, logConfiguration.Path)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	routinesPool := safe.NewPool(ctx, 10000)

	// backend provider
	backend := backendProvider.NewBackend()
	fp := fileProvider.NewFile()
	w := watcher.NewConfigurationWatcher(routinesPool)
	w.AddProvider(backend, fp)

	upstreamFactory := upstream.NewFactory(config.DefaultConfig.Gateway, routinesPool)

	routerFactory := routerManager.NewRouterFactory(config.DefaultConfig.Gateway, upstreamFactory)

	// server start
	server := gateway.NewServer(routinesPool, w, routerFactory, &config.DefaultConfig.Gateway)

	server.Start(ctx)
	defer server.Close()

	server.Wait()

}
