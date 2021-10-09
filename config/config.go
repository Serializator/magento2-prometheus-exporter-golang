package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Magento struct {
		Url    string
		Bearer string
	}
}

func NewConfig() (*Config, error) {
	config := &Config{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error whilst reading config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("fatal error whilst unmarshaling config file: %w", err)
	}

	return config, nil
}
