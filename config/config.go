package config

import (
	"github.com/spf13/viper"
)

// IConfig is an interface of config
type IConfig interface {
	InitConfig() error
	GetConfig() Config
}

// Config file
type Config struct {
	InteractiveFilterCommand string
}

// InitConfig initializes config
func (c Config) InitConfig() (err error) {
	viper.SetConfigName("gecco")
	viper.AddConfigPath("$HOME/.config")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	// if err = viper.Unmarshal(&config.Conf); err != nil {
	if err = viper.Unmarshal(&c); err != nil {
		return
	}

	return
}

// GetConfig returns config
func (c Config) GetConfig() Config {
	return c
}
