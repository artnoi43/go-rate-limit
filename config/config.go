package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MaxGuard int    `mapstructure:"max_guard"`
	URL      string `mapstructure:"url"`
}

func Load() (*Config, error) {
	viper.SetDefault("max_guard", 10)
	viper.SetDefault("url", "https://google.com")

	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		// Config file not found
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// WriteConfig() just won't create new file if doesn't exist
			viper.SafeWriteConfig()
		} else {
			return nil, err
		}
	}

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, err
}
