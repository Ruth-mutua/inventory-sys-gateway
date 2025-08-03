package config

import (
	"os"

	"go-inventory-system/shared"

	"gopkg.in/yaml.v3"
)

// GatewayConfig holds gateway-specific configuration
type GatewayConfig struct {
	Port   string         `yaml:"port"`
	Routes []shared.Route `yaml:"routes"`
}

// LoadConfig loads gateway configuration
func LoadConfig() *GatewayConfig {
	return &GatewayConfig{
		Port: getEnv("GATEWAY_PORT", "8000"),
	}
}

// LoadRoutes loads routes from YAML file
func LoadRoutes(filename string) ([]shared.Route, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config struct {
		Routes []shared.Route `yaml:"routes"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Routes, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
