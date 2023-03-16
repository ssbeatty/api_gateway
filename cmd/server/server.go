package main

import (
	backendServer "api_gateway/internal/backend"
	backendConfig "api_gateway/internal/backend/config"
	"api_gateway/internal/backend/handler"
	"api_gateway/internal/backend/models"
	"api_gateway/internal/gateway"
	"api_gateway/internal/gateway/config"
	routerManager "api_gateway/internal/gateway/manager/router"
	"api_gateway/internal/gateway/manager/upstream"
	backendProvider "api_gateway/internal/gateway/provider/backend"
	"api_gateway/internal/gateway/watcher"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/safe"
	"context"
	"os/signal"
	"syscall"
)

func main() {
	// init default config
	config.SetupConfig()
	backendConfig.SetupConfig()

	// setup log
	logConfiguration := config.DefaultConfig.Log
	logs.SetupLogger(logConfiguration.Level, logConfiguration.Path)

	//system signal
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	routinesPool := safe.NewPool(ctx, -1)

	// init backend database
	backendCfg := backendConfig.DefaultConfig
	err := models.InitModels(backendCfg.DB, ctx)
	if err != nil {
		panic(err)
	}
	// only backend provider
	backend := backendProvider.NewBackend()
	webConfig := backendCfg.WebServer

	backendService := handler.NewService(webConfig, backend)
	backendServer.Serve(backendService)

	w := watcher.NewConfigurationWatcher(routinesPool)
	w.AddProvider(backend)

	upstreamFactory := upstream.NewFactory(config.DefaultConfig.Gateway, routinesPool)

	routerFactory := routerManager.NewRouterFactory(config.DefaultConfig.Gateway, upstreamFactory)

	// server start
	server := gateway.NewServer(routinesPool, w, routerFactory, &config.DefaultConfig.Gateway)

	server.Start(ctx)
	defer server.Close()

	server.Wait()

}
