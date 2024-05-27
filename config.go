package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Define a struct that matches the structure of your YAML file
type Config struct {
	EWeChat       EWeChat `mapstructure:"ewechat"`
	DiskUsageRate float64 `mapstructure:"disk_usage_rate"`
	CpuUsageRate  float64 `mapstructure:"cpu_usage_rate"`
	MemUsageRate  float64 `mapstructure:"mem_usage_rate"`
}

type EWeChat struct {
	CorpID     string `mapstructure:"corp_id"`
	CorpSecret string `mapstructure:"corp_secret"`
	AgentID    int    `mapstructure:"agent_id"`
	Receivers  string `mapstructure:"receivers"`
}

func GetConfig() Config {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	// Resolve the directory of the executable.
	exePath := filepath.Dir(exe)

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // path to look for the config file in
	viper.AddConfigPath(exePath)

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %w", err))
	}

	return config
}

var CFG = GetConfig()
