package file

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/safe"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"reflect"
)

// NewFile new file config provider.
// Power by spf13/viper, notify config when file change
func NewFile() *File {
	v := viper.New()

	return &File{
		applyMessage: make(chan dynamic.Message, 10),
		vipHandler:   v,
		endpointsCfs: make([]config.Endpoint, 0),
	}
}

type Config struct {
	Endpoints []config.Endpoint `yaml:"endpoints"`
}

// File file config provider.
type File struct {
	applyMessage chan dynamic.Message
	vipHandler   *viper.Viper
	endpointsCfs []config.Endpoint
	logger       zerolog.Logger
}

func (b *File) setupConfigs() error {
	b.vipHandler.SetConfigName("endpoints")
	b.vipHandler.SetConfigType("yaml")
	b.vipHandler.AddConfigPath("/etc/api_gateway/")
	b.vipHandler.AddConfigPath("$HOME/.api_gateway")
	b.vipHandler.AddConfigPath(".")

	err := b.vipHandler.ReadInConfig()
	if err != nil {
		return err
	}

	b.reloadConfigs()

	b.vipHandler.OnConfigChange(func(e fsnotify.Event) {
		b.logger.Debug().Str("Config file changed:", e.Name)
		b.reloadConfigs()

	})
	b.vipHandler.WatchConfig()

	return nil
}

func (b *File) reloadConfigs() {
	endpointCfg := Config{}

	err := b.vipHandler.Unmarshal(&endpointCfg)
	if err != nil {
		b.logger.Error().Err(err).Msg("Error when marshal endpoints config")
	}
	if reflect.DeepEqual(b.endpointsCfs, endpointCfg.Endpoints) {
		return
	}

	b.applyMessage <- dynamic.Message{
		ProviderName:  b.Name(),
		Configuration: endpointCfg.Endpoints,
	}

	b.endpointsCfs = endpointCfg.Endpoints

}

func (b *File) Provide(configurationChan chan<- dynamic.Message, pool *safe.Pool) error {
	pool.Go(func() {
		for msg := range b.applyMessage {
			configurationChan <- msg
		}
	})
	return nil
}

func (b *File) Init() error {
	b.logger = log.With().Str(logs.ProviderName, b.Name()).Logger()

	err := b.setupConfigs()
	if err != nil {
		return err
	}

	return nil
}

func (b *File) Name() string {
	return "file"
}
