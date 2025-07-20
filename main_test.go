package main

import (
	"osc2hue/internal/config"
	"testing"
)

func TestConfigStructure(t *testing.T) {
	// Test that our config structure is valid
	cfg := &config.Config{
		OSC: config.OSCConfig{
			Host: "127.0.0.1",
			Port: 9000,
		},
		Hue: config.HueConfig{
			BridgeIP: "192.168.1.100",
			APIKey:   "test-api-key",
		},
	}

	if cfg.OSC.Host != "127.0.0.1" {
		t.Errorf("Expected OSC host to be 127.0.0.1, got %s", cfg.OSC.Host)
	}

	if cfg.OSC.Port != 9000 {
		t.Errorf("Expected OSC port to be 9000, got %d", cfg.OSC.Port)
	}

	if cfg.Hue.BridgeIP != "192.168.1.100" {
		t.Errorf("Expected Hue bridge IP to be 192.168.1.100, got %s", cfg.Hue.BridgeIP)
	}

	if cfg.Hue.APIKey != "test-api-key" {
		t.Errorf("Expected Hue API key to be test-api-key, got %s", cfg.Hue.APIKey)
	}
}
