package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// App
	Port              string        `mapstructure:"APP_PORT"`
	TokenSymmetricKey string        `mapstructure:"APP_TOKEN_SYMMETRIC_KEY"`
	TokenExpiration   time.Duration `mapstructure:"APP_TOKEN_EXPIRATION"`

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
