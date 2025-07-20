package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GEMINI_API_KEY string `yaml:"gemini_api_key"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if configPath == "" {
		configPath = "config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("error parsing config file: %v", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		cfg.GEMINI_API_KEY = apiKey
	}

	return cfg, nil
}
