package hue

import (
	"fmt"
	"time"

	"github.com/openhue/openhue-go"
)

// Bridge represents a discovered Hue bridge
type Bridge struct {
	IPAddress string
}

// DiscoverBridge discovers Hue bridges on the network
func DiscoverBridge(timeout time.Duration) (*Bridge, error) {
	// Create a discovery client with timeout
	discovery := openhue.NewBridgeDiscovery(openhue.WithTimeout(timeout))

	// Discover bridges
	bridge, err := discovery.Discover()
	if err != nil {
		return nil, fmt.Errorf("bridge discovery failed: %v", err)
	}

	return &Bridge{
		IPAddress: bridge.IpAddress,
	}, nil
}

// AuthenticateWithBridge performs bridge authentication
func AuthenticateWithBridge(bridgeIP string) (string, error) {
	if bridgeIP == "" {
		return "", fmt.Errorf("bridge IP not set")
	}

	// Create authenticator
	authenticator, err := openhue.NewAuthenticator(bridgeIP)
	if err != nil {
		return "", fmt.Errorf("failed to create authenticator: %v", err)
	}

	// Keep trying to authenticate until button is pressed or we get an error
	var apiKey string
	for len(apiKey) == 0 {
		key, retry, err := authenticator.Authenticate()

		if err != nil && retry {
			// Link button not pressed yet, continue waiting
			time.Sleep(500 * time.Millisecond)
		} else if err != nil && !retry {
			// Real error occurred
			return "", fmt.Errorf("authentication failed: %v", err)
		} else {
			// Success!
			apiKey = key
		}
	}

	return apiKey, nil
}

// IsValidAPIKey checks if an API key is valid (not empty or placeholder)
func IsValidAPIKey(apiKey string) bool {
	return apiKey != "" && apiKey != "your-hue-api-key-here"
}
