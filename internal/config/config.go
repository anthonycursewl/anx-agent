package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	GEMINI_API_KEY string `yaml:"gemini_api_key"`
}

// LoadConfig loads the configuration from the specified file path
// If configPath is empty, it will try to load from environment variables
func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	// First try to load from config file if path is provided
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("error parsing config file: %v", err)
		}
	}

	// Override with environment variables if they exist
	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		cfg.GEMINI_API_KEY = apiKey
	}

	return cfg, nil
}
