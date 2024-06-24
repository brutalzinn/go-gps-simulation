package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	GoogleMapsAPIKey string `yaml:"google_maps_api_key"`
	DefaultDevice    string `yaml:"default_device"`
	Port             int    `yaml:"port"`
	AdbBaseURL       string `yaml:"adb_base_url"`
	AdbPath          string `yaml:"adb_path"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config file: %w", err)
	}

	return &config, nil
}
