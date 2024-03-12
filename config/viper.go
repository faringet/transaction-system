package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

func NewViper(filename string) (*viper.Viper, *Config, error) {
	v := viper.New()

	if filename != "" {
		v.SetConfigName(filename)
		v.AddConfigPath(".")
		v.AddConfigPath(filepath.FromSlash("config"))
	}

	err := v.ReadInConfig()
	if err != nil {
		log.Println("viper failed to read app config file:", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return v, &config, nil
}
