package main

import (
	"api_gateway/internal/gateway"
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/internal/gateway/provider"
	"api_gateway/internal/gateway/watcher"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/safe"
	"context"
	"os/signal"
	"syscall"
	"time"
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

	// todo test code
	go func() {
		time.Sleep(time.Second * 2)

		backend.ReloadConfig(dynamic.Message{
			ProviderName: backend.Name(),
			Configuration: map[string]config.Endpoint{
				"test": {
					Name:       "test",
					ListenPort: 8080,
				},
			},
		})
	}()

	// server start
	server := gateway.NewServer(routinesPool, w)

	server.Start(ctx)
	defer server.Close()

	server.Wait()

}
