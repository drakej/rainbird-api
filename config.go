package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig loads the configuration file (config.toml or config.yaml)
func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}
