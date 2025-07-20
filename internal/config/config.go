package config

import (
	"encoding/json"
	"os"
)

// Config holds the application configuration
type Config struct {
	OSC OSCConfig `json:"osc"`
	Hue HueConfig `json:"hue"`
}

// OSCConfig holds OSC server configuration
type OSCConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

// HueConfig holds Philips Hue configuration
type HueConfig struct {
	BridgeIP string `json:"bridge_ip"`
	APIKey   string `json:"api_key"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(config *Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}
