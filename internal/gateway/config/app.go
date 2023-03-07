package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var (
	DefaultConfig = APP{}
)

func SetupConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/api_gateway/")
	viper.AddConfigPath("$HOME/.api_gateway")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&DefaultConfig)
	if err != nil {
		panic(err)
	}
}

// APP App logs
type APP struct {
	Log     Log     `yaml:"log"`
	Gateway Gateway `yaml:"gateway"`
}

type Log struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}

type Gateway struct {
	GraceTimeOut     time.Duration `yaml:"graceTimeOut"`
	ListenUDPTimeOut time.Duration `yaml:"listenUDPTimeOut"`
}
