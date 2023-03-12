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
	Log     Log     `yaml:"log" mapstructure:"log"`
	Gateway Gateway `yaml:"gateway" mapstructure:"gateway"`
}

func NewAPP() APP {
	return APP{
		Log:     newLog(),
		Gateway: newGateway(),
	}
}

type Log struct {
	Level     string `yaml:"level" mapstructure:"level"`
	Path      string `yaml:"path" mapstructure:"path"`
	AccessLog `yaml:"access_log" mapstructure:"access_log"`
}

type AccessLog struct {
	Enable      bool   `yaml:"enable" mapstructure:"enable"`
	HttpLogPath string `yaml:"httpLogPath" mapstructure:"httpLogPath"`
}

func newLog() Log {
	return Log{
		Level:     "info",
		Path:      "stdout",
		AccessLog: newAccessLog(),
	}
}

func newAccessLog() AccessLog {
	return AccessLog{
		Enable:      false,
		HttpLogPath: "stdout",
	}
}

type Gateway struct {
	GraceTimeOut     time.Duration `yaml:"graceTimeOut" mapstructure:"graceTimeOut"`
	ListenUDPTimeOut time.Duration `yaml:"listenUDPTimeOut" mapstructure:"listenUDPTimeOut"`
	HTTPReadTimeOut  time.Duration `yaml:"httpReadTimeout" mapstructure:"httpReadTimeout"`
	HTTPWriteTimeOut time.Duration `yaml:"httpWriteTimeout" mapstructure:"httpWriteTimeout"`
	HTTPIdleTimeOut  time.Duration `yaml:"httpIdleTimeout" mapstructure:"httpIdleTimeout"`
}

func newGateway() Gateway {
	return Gateway{
		GraceTimeOut:     time.Nanosecond * 50,
		ListenUDPTimeOut: time.Second * 10,
		HTTPReadTimeOut:  time.Second * 30,
		HTTPWriteTimeOut: time.Second * 30,
		HTTPIdleTimeOut:  time.Second * 120,
	}
}
