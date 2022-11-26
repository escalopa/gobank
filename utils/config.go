package utils

import (
	"github.com/spf13/viper"
)

type DBConfig struct {
	ConnectionString string `mapstructure:"connection_string"`
	Driver           string `mapstructure:"driver"`
}

type AppConfig struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	App AppConfig `mapstructure:"app"`
	DB  DBConfig  `mapstructure:"database"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
