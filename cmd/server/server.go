package main

import (
	"api_gateway/internal/gateway"
	"api_gateway/internal/gateway/provider"
	"api_gateway/internal/gateway/watcher"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/safe"
	"context"
	"os/signal"
	"syscall"
)

func init() {
	logs.SetupLogger("debug", "stdout")
}

func main() {
	// init default
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	routinesPool := safe.NewPool(ctx, 10000)

	// backend provider
	backend := provider.NewBackend()
	w := watcher.NewConfigurationWatcher(routinesPool)

	w.AddProvider(backend)

	// server start
	server := gateway.NewServer(routinesPool, w)

	server.Start(ctx)
	defer server.Close()

	server.Wait()

}
