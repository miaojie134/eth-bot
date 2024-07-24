package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AlpacaAPIKey    string
	AlpacaAPISecret string
	DBPath          string
	ServerPort      string
	LogDir          string
	LogLevel        string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	err := viper.Unmarshal(&config)
	return &config, err
}
