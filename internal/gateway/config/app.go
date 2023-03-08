package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var (
	DefaultConfig = NewAPP()
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

func NewAPP() APP {
	return APP{
		Log:     newLog(),
		Gateway: newGateway(),
	}
}

type Log struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}

func newLog() Log {
	return Log{
		Level: "info",
		Path:  "stdout",
	}
}

type Gateway struct {
	GraceTimeOut     time.Duration `yaml:"graceTimeOut"`
	ListenUDPTimeOut time.Duration `yaml:"listenUDPTimeOut"`
	HTTPReadTimeOut  time.Duration `yaml:"httpReadTimeout"`
	HTTPWriteTimeOut time.Duration `yaml:"httpWriteTimeout"`
	HTTPIdleTimeOut  time.Duration `yaml:"httpIdleTimeout"`
}

func newGateway() Gateway {
	return Gateway{
		GraceTimeOut:     time.Nanosecond * 50,
		ListenUDPTimeOut: time.Second * 30,
		HTTPReadTimeOut:  time.Second * 30,
		HTTPWriteTimeOut: time.Second * 30,
		HTTPIdleTimeOut:  time.Second * 120,
	}
}
