package config

import (
	"fmt"
	"github.com/spf13/viper"
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
	Log       Log       `yaml:"log" mapstructure:"log"`
	DB        DB        `yaml:"db" mapstructure:"db"`
	WebServer WebServer `yaml:"web" mapstructure:"web"`
}

func NewAPP() APP {
	return APP{
		Log:       newLog(),
		DB:        newDB(),
		WebServer: newWebServer(),
	}
}

type Log struct {
	Level string `yaml:"level" mapstructure:"level"`
	Path  string `yaml:"path" mapstructure:"path"`
}

type AccessLog struct {
	Enable      bool   `yaml:"enable" mapstructure:"enable"`
	HttpLogPath string `yaml:"httpLogPath" mapstructure:"httpLogPath"`
}

func newLog() Log {
	return Log{
		Level: "info",
		Path:  "stdout",
	}
}

type DB struct {
	Driver   string `yaml:"driver" mapstructure:"driver"`
	DSN      string `yaml:"dsn" mapstructure:"dsn"`
	User     string `yaml:"user" mapstructure:"user"`
	Pass     string `yaml:"pass" mapstructure:"pass"`
	DBName   string `yaml:"db_name" mapstructure:"db_name"`
	DataPath string `yaml:"data_path" mapstructure:"data_path"`
}

func newDB() DB {
	return DB{
		Driver:   "sqlite",
		DSN:      "127.0.0.1:3306",
		User:     "root",
		Pass:     "123456",
		DBName:   "gateway",
		DataPath: "db",
	}
}

type WebServer struct {
	Name string `yaml:"name" mapstructure:"name"`
	Addr string `yaml:"addr" mapstructure:"addr"`
	Port int    `yaml:"port" mapstructure:"port"`
}

func newWebServer() WebServer {
	return WebServer{
		Name: "api_gateway",
		Addr: "0.0.0.0",
		Port: 8099,
	}
}
