package hue

import (
	"testing"
)

func TestIsValidAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		expected bool
	}{
		{
			name:     "Empty API key",
			apiKey:   "",
			expected: false,
		},
		{
			name:     "Default placeholder",
			apiKey:   "your-hue-api-key-here",
			expected: false,
		},
		{
			name:     "Valid API key",
			apiKey:   "abcd1234efgh5678ijkl9012mnop3456qrst7890",
			expected: true,
		},
		{
			name:     "Short API key",
			apiKey:   "short",
			expected: true, // Any non-empty, non-placeholder key is considered valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidAPIKey(tt.apiKey)
			if result != tt.expected {
				t.Errorf("IsValidAPIKey(%q) = %v, expected %v", tt.apiKey, result, tt.expected)
			}
		})
	}
}

func TestBridge(t *testing.T) {
	// Test Bridge structure
	bridge := &Bridge{
		IPAddress: "192.168.1.100",
	}

	if bridge.IPAddress != "192.168.1.100" {
		t.Errorf("Expected bridge IP 192.168.1.100, got %s", bridge.IPAddress)
	}
}
