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
		currentConfigs:         make(map[string]dynamic.Message),
		mu:                     &sync.Mutex{},
		routinesPool:           pool,
	}
}

// ConfigurationWatcher watches configuration changes.
type ConfigurationWatcher struct {
	providers              []Provider
	allProvidersConfigs    chan dynamic.Message
	configurationListeners []func(dynamic.Configuration)
	currentConfigs         map[string]dynamic.Message

	mu           *sync.Mutex
	routinesPool *safe.Pool
}

// Start the configuration watcher.
func (c *ConfigurationWatcher) Start() {
	c.routinesPool.GoCtx(c.receiveConfigurations)
	c.startProviderAggregator()
}

func (c *ConfigurationWatcher) AddProvider(p ...Provider) {
	if c.providers == nil {
		c.providers = make([]Provider, 0)
	}

	c.providers = append(c.providers, p...)
}

// AddListener adds a new listener function used when new configuration is provided.
func (c *ConfigurationWatcher) AddListener(listener func(dynamic.Configuration)) {
	if c.configurationListeners == nil {
		c.configurationListeners = make([]func(dynamic.Configuration), 0)
	}
	c.configurationListeners = append(c.configurationListeners, listener)
}

func (c *ConfigurationWatcher) startProviderAggregator() {
	for idx := range c.providers {
		provider := c.providers[idx]
		c.routinesPool.Go(func() {
			if err := provider.Init(); err != nil {
				log.Error().Err(err).Msgf("Error when init provider %s", provider.Name())
				return
			}
			log.Info().Msgf("Starting provider %s", provider.Name())
			err := provider.Provide(c.allProvidersConfigs, c.routinesPool)
			if err != nil {
				log.Error().Err(err).Msgf("Error starting provider: %s", provider.Name())
			}
		})
	}
}

func (c *ConfigurationWatcher) applyConfigurations(conf dynamic.Configuration) {
	log.Debug().Interface("applyConfig", conf).Msg("Apply New Configuration")

	for _, listener := range c.configurationListeners {
		listener(conf)
	}
}

func isExistedName(endpoints []config.Endpoint, name string) (config.Endpoint, bool) {
	for _, endpoint := range endpoints {
		if name == endpoint.Name {
			return endpoint, true
		}
	}

	return config.Endpoint{}, false
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

			currentConfigs := c.currentConfigs[newConfig.ProviderName]

			if reflect.DeepEqual(currentConfigs, newConfig) {
				continue
			}
			c.mu.Lock()

			log.Info().Msgf("Change new config from Provider: %s", newConfig.ProviderName)

			newConfigs := newConfig.Configuration
			for _, currentConfig := range currentConfigs.Configuration {
				if newOne, existed := isExistedName(newConfigs, currentConfig.Name); !existed {
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
			for _, newOne := range newConfigs {
				if _, existed := isExistedName(currentConfigs.Configuration, newOne.Name); !existed {
					c.applyConfigurations(
						dynamic.Configuration{
							Action:   ActionCreate,
							EndPoint: newOne,
						})
				}
			}

			c.currentConfigs[newConfig.ProviderName] = newConfig

			c.mu.Unlock()
		}
	}
}
