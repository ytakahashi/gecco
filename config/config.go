package config

import (
	"github.com/spf13/viper"
)

// IConfig is an interface of config
type IConfig interface {
	InitConfig() error
	GetConfig() Config
}

// Config stores values read from config file.
type Config struct {
	// InteractiveFilterCommand is used when gecco is run with "-i" option.
	// This holds interactive filter command (like "fzf" ot "peco") as string.
	InteractiveFilterCommand string
}

// InitConfig initializes config object.
// Config file should be placed at "~/.config/" directory.
func (c *Config) InitConfig() (err error) {
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
