package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Database       DatabaseConfig `yaml:"database"`
	ProviderConfig ProviderConfig `yaml:"provider"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string `yaml:"driver"` // e.g., "postgres"
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslmode"`
}

type ProviderConfig struct {
	MockOne MockOneConfig `yaml:"mock-one"`
	MockTwo MockTwoConfig `yaml:"mock-two"`
}

type MockOneConfig struct {
	Url string `yaml:"url"`
}

type MockTwoConfig struct {
	Url string `yaml:"url"`
}

var config *Config

// Load reads the configuration file and returns a Config struct
func Load() (*Config, error) {
	if config == nil {
		data, err := os.ReadFile("config.yaml")
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("error parsing config file: %w", err)
		}

		config = &cfg
	}

	return config, nil
}
