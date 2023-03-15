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

func newLog() Log {
	return Log{
		Level: "info",
		Path:  "stdout",
	}
}

type DB struct {
	Driver   string `yaml:"driver" mapstructure:"driver"`
	DSN      string `yaml:"dsn" mapstructure:"dsn"`
	UserName string `yaml:"username" mapstructure:"username"`
	PassWord string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"db_name" mapstructure:"db_name"`
	DataPath string `yaml:"data_path" mapstructure:"data_path"`
}

func newDB() DB {
	return DB{
		Driver:   "sqlite",
		DSN:      "127.0.0.1:3306",
		UserName: "root",
		PassWord: "123456",
		DBName:   "gateway",
		DataPath: "data",
	}
}

type WebServer struct {
	BindAddr string `yaml:"bind_addr" mapstructure:"bind_addr"`
	BindPort int    `yaml:"bind_port" mapstructure:"bind_port"`
	Jwt      Jwt    `yaml:"jwt" mapstructure:"jwt"`
}

func newWebServer() WebServer {
	return WebServer{
		BindAddr: "0.0.0.0",
		BindPort: 8099,
		Jwt:      newJwt(),
	}
}

type Jwt struct {
	BearerSchema  string        `yaml:"bearer_schema" mapstructure:"bearer_schema"`
	JwtSecretPath string        `yaml:"jwt_secret" mapstructure:"jwt_secret"`
	JwtExp        time.Duration `yaml:"jwt_exp" mapstructure:"jwt_exp"`
}

func newJwt() Jwt {
	return Jwt{
		BearerSchema:  "Bearer",
		JwtSecretPath: "00163e06360c",
		JwtExp:        2 * time.Hour,
	}
}
