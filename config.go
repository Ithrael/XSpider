package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Headers     map[string]string `mapstructure:"Headers"`
	Restriction struct {
		MaxDepth           int      `mapstructure:"MaxDepth"`
		MaxCount           int      `mapstructure:"MaxCount"`
		AllowedDomains     []string `mapstructure:"AllowedDomains"`
		ExcludedDomains    []string `mapstructure:"ExcludedDomains"`
		AllowedPaths       []string `mapstructure:"AllowedPaths"`
		ExcludedPaths      []string `mapstructure:"ExcludedPaths"`
		AllowedQueryKey    []string `mapstructure:"AllowedQueryKey"`
		ExcludedQueryKey   []string `mapstructure:"ExcludedQueryKey"`
		Parallelism        int      `mapstructure:"Parallelism"`
		RandomDelayMaxTime int      `mapstructure:"RandomDelayMaxTime"`
	} `mapstructure:"Restriction"`
}

// Function to load the configuration from the YAML file
func LoadConfig() (*Config, error) {
	var config Config
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode config struct, Error is %s", err)
		return nil, err
	}
	return &config, nil
}
