package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// App
	HTTPPort               string        `mapstructure:"HTTP_PORT"`
	GRPCPort               string        `mapstructure:"GRPC_PORT"`
	TokenSymmetricKey      string        `mapstructure:"APP_TOKEN_SYMMETRIC_KEY"`
	AccessTokenExpiration  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRATION"`
	RefreshTokenExpiration time.Duration `mapstructure:"REFRESH__TOKEN_EXPIRATION"`

	// Database
	ConnectionString string `mapstructure:"DB_CONNECTION_STRING"`
	Driver           string `mapstructure:"DB_DRIVER"`
	Name             string `mapstructure:"DB_NAME"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app.env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
