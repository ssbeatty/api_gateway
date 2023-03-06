package watcher

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/pkg/safe"
	"context"
	"github.com/rs/zerolog/log"
	"reflect"
	"sync"
)

const (
	ActionDelete = "delete"
	ActionUpdate = "update"
	ActionCreate = "create"
)

func NewConfigurationWatcher(pool *safe.Pool) *ConfigurationWatcher {
	return &ConfigurationWatcher{
		allProvidersConfigs:    make(chan dynamic.Message, 100),
		configurationListeners: make([]func(dynamic.Configuration), 0),
		currentConfigs:         make(map[string]config.Endpoint),
		mu:                     &sync.Mutex{},
		routinesPool:           pool,
	}
}

// ConfigurationWatcher watches configuration changes.
type ConfigurationWatcher struct {
	providers              []Provider
	allProvidersConfigs    chan dynamic.Message
	configurationListeners []func(dynamic.Configuration)
	currentConfigs         map[string]config.Endpoint

	mu           *sync.Mutex
	routinesPool *safe.Pool
}

// Start the configuration watcher.
func (c *ConfigurationWatcher) Start() {
	c.routinesPool.GoCtx(c.receiveConfigurations)
	c.startProviderAggregator()
}

func (c *ConfigurationWatcher) AddProvider(p Provider) {
	if c.providers == nil {
		c.providers = make([]Provider, 0)
	}

	c.providers = append(c.providers, p)
}

// AddListener adds a new listener function used when new configuration is provided.
func (c *ConfigurationWatcher) AddListener(listener func(dynamic.Configuration)) {
	if c.configurationListeners == nil {
		c.configurationListeners = make([]func(dynamic.Configuration), 0)
	}
	c.configurationListeners = append(c.configurationListeners, listener)
}

func (c *ConfigurationWatcher) startProviderAggregator() {
	for _, provider := range c.providers {
		c.routinesPool.Go(func() {
			log.Info().Msgf("Starting provider %s", provider.Name())
			err := provider.Provide(c.allProvidersConfigs, c.routinesPool)
			if err != nil {
				log.Error().Err(err).Msgf("Error starting provider: %s", provider.Name())
			}
		})
	}
}

func (c *ConfigurationWatcher) applyConfigurations(conf dynamic.Configuration) {
	for _, listener := range c.configurationListeners {
		listener(conf)
	}
}

func (c *ConfigurationWatcher) receiveConfigurations(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case newConfig, ok := <-c.allProvidersConfigs:
			if !ok {
				return
			}

			if reflect.DeepEqual(c.currentConfigs, newConfig.Configuration) {
				continue
			}
			c.mu.Lock()

			log.Info().Msgf("Change new config from Provider: %s", newConfig.ProviderName)

			newConfigs := newConfig.Configuration
			for name, currentConfig := range c.currentConfigs {
				if newOne, existed := newConfigs[name]; !existed {
					c.applyConfigurations(
						dynamic.Configuration{
							Action:   ActionDelete,
							EndPoint: currentConfig,
						})
				} else {
					if reflect.DeepEqual(currentConfig, newOne) {
						continue
					}
					c.applyConfigurations(
						dynamic.Configuration{
							Action:   ActionUpdate,
							EndPoint: newOne,
						})
				}
			}
			for name, newOne := range newConfigs {
				if _, existed := c.currentConfigs[name]; !existed {
					c.applyConfigurations(
						dynamic.Configuration{
							Action:   ActionCreate,
							EndPoint: newOne,
						})
				}
			}

			c.currentConfigs = newConfig.Configuration

			c.mu.Unlock()
		}
	}
}
