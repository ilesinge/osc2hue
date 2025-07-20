package osc

import (
	"testing"
)

func TestNewServer(t *testing.T) {
	// Test creating a new OSC server
	server := NewServer("127.0.0.1", 8080)

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}
}

func TestServerConfiguration(t *testing.T) {
	// Test that server can be created with different configurations
	tests := []struct {
		name string
		host string
		port int
	}{
		{
			name: "Localhost",
			host: "127.0.0.1",
			port: 8080,
		},
		{
			name: "All interfaces",
			host: "0.0.0.0",
			port: 9000,
		},
		{
			name: "Specific IP",
			host: "192.168.1.100",
			port: 5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.host, tt.port)
			if server == nil {
				t.Errorf("Failed to create server with host %s and port %d", tt.host, tt.port)
			}
		})
	}
}
