package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test loading non-existent file
	_, err := LoadConfig("nonexistent.json")
	if err == nil {
		t.Error("Expected error when loading non-existent config file")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.json")

	// Create test config
	originalConfig := &Config{
		OSC: OSCConfig{
			Host: "127.0.0.1",
			Port: 9000,
		},
		Hue: HueConfig{
			BridgeIP: "192.168.1.100",
			APIKey:   "test-api-key-123",
		},
	}

	// Save config
	err := SaveConfig(originalConfig, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load config back
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify config values
	if loadedConfig.OSC.Host != originalConfig.OSC.Host {
		t.Errorf("Expected OSC host %s, got %s", originalConfig.OSC.Host, loadedConfig.OSC.Host)
	}

	if loadedConfig.OSC.Port != originalConfig.OSC.Port {
		t.Errorf("Expected OSC port %d, got %d", originalConfig.OSC.Port, loadedConfig.OSC.Port)
	}

	if loadedConfig.Hue.BridgeIP != originalConfig.Hue.BridgeIP {
		t.Errorf("Expected Hue bridge IP %s, got %s", originalConfig.Hue.BridgeIP, loadedConfig.Hue.BridgeIP)
	}

	if loadedConfig.Hue.APIKey != originalConfig.Hue.APIKey {
		t.Errorf("Expected Hue API key %s, got %s", originalConfig.Hue.APIKey, loadedConfig.Hue.APIKey)
	}
}
