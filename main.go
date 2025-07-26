package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"osc2hue/internal/config"
	"osc2hue/internal/hue"
	"osc2hue/internal/osc"

	"github.com/openhue/openhue-go"
)

func main() {
	configPath := "config.json"

	// Load and setup configuration
	cfg := loadOrCreateConfig(configPath)

	// Setup bridge discovery and authentication
	setupBridgeConnection(cfg, configPath)

	// Create client and discover lights
	home, lights := setupHueClient(cfg)

	// Setup and start OSC server
	startOSCServer(cfg, home, lights)
}

// startOSCServer creates, configures and starts the OSC server
func startOSCServer(cfg *config.Config, home *openhue.Home, lights []openhue.LightGet) {
	// Create OSC server
	oscServer := osc.NewServer(cfg.OSC.Host, cfg.OSC.Port)

	// Add all OSC handlers
	addAllHandlers(oscServer, home, lights)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down...")
		oscServer.Stop()
		os.Exit(0)
	}()

	// Start the OSC server
	log.Printf("Starting OSC2Hue bridge...")
	log.Printf("OSC Server: %s:%d", cfg.OSC.Host, cfg.OSC.Port)
	log.Printf("Hue Bridge: %s", cfg.Hue.BridgeIP)
	log.Printf("Available OSC commands:")
	log.Printf("  /hue/{id}/on {0|1} [duration_ms]")
	log.Printf("  /hue/{id}/set {x|-1} [y|-1] [brightness|-1] [duration_ms|-1]")
	log.Printf("  /hue/{id}/brightness {0-254} [duration_ms]")
	log.Printf("  /hue/{id}/color {x} {y} [duration_ms]")
	log.Printf("  /hue/all/on {0|1} [duration_ms]")
	log.Printf("  /hue/all/set {x|-1} [y|-1] [brightness|-1] [duration_ms|-1]")
	log.Printf("  /hue/all/brightness {0-254} [duration_ms]")
	log.Printf("  /hue/all/color {x} {y} [duration_ms]")
	log.Printf("Note: Use -1 for null values in /set commands to skip color, brightness, or duration")

	if err := oscServer.Start(); err != nil {
		log.Fatalf("Failed to start OSC server: %v", err)
	}
}

// setupHueClient creates the Hue client and discovers lights
func setupHueClient(cfg *config.Config) (*openhue.Home, []openhue.LightGet) {
	// Create client for Hue API
	home, err := openhue.NewHome(cfg.Hue.BridgeIP, cfg.Hue.APIKey)
	if err != nil {
		log.Printf("Failed to create Hue client: %v", err)
		log.Printf("Continuing anyway - you can test OSC messages but they won't control lights")
		return nil, nil
	}

	var lights []openhue.LightGet
	if home != nil {
		// Test connection and discover lights
		log.Printf("Testing connection to Hue Bridge at %s...", cfg.Hue.BridgeIP)
		lightsMap, err := home.GetLights()
		if err != nil {
			log.Printf("Warning: Failed to connect to Hue Bridge: %v", err)
			log.Printf("Continuing anyway - you can test OSC messages but they won't control lights")
			return home, lights
		} else {
			// Convert map to slice
			for _, light := range lightsMap {
				lights = append(lights, light)
			}
			log.Printf("Successfully connected! Found %d lights:", len(lights))
			for id, light := range lights {
				log.Printf("  Light #%d %s: %s", id+1, *light.Id, *light.Metadata.Name)
			}
		}
	}
	return home, lights
}

// loadOrCreateConfig loads configuration from file or creates a default one
func loadOrCreateConfig(configPath string) *config.Config {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		log.Println("Using example configuration...")
		cfg = &config.Config{
			OSC: config.OSCConfig{
				Host: "0.0.0.0",
				Port: 8080,
			},
			Hue: config.HueConfig{
				BridgeIP: "", // Empty, will be discovered
				APIKey:   "", // Empty, will be authenticated
			},
		}
	}
	return cfg
}

// setupBridgeConnection handles bridge discovery and authentication
func setupBridgeConnection(cfg *config.Config, configPath string) {
	// Discover bridge if IP is not set or seems invalid
	if cfg.Hue.BridgeIP == "" {
		discoverAndSaveBridge(cfg, configPath)
	}

	// Authenticate with bridge if needed
	if !hue.IsValidAPIKey(cfg.Hue.APIKey) {
		authenticateAndSaveAPIKey(cfg, configPath)
	} else {
		log.Printf("Using existing API key: %s", cfg.Hue.APIKey)
	}
}

// discoverAndSaveBridge discovers the Hue bridge and saves the IP to config
func discoverAndSaveBridge(cfg *config.Config, configPath string) {
	log.Println("Discovering Hue bridges...")
	bridge, err := hue.DiscoverBridge(5 * time.Second)
	if err != nil {
		log.Printf("Bridge discovery failed: %v", err)
		log.Println("Please manually set the bridge_ip in config.json")
		return
	}

	log.Printf("Found Hue bridge at %s", bridge.IPAddress)

	// Update config if the bridge IP has changed or was empty
	if cfg.Hue.BridgeIP != bridge.IPAddress {
		cfg.Hue.BridgeIP = bridge.IPAddress
		log.Printf("Updated bridge IP to %s", bridge.IPAddress)

		// Save the updated configuration
		if err := config.SaveConfig(cfg, configPath); err != nil {
			log.Printf("Warning: Failed to save updated config: %v", err)
		} else {
			log.Println("Configuration saved with discovered bridge IP")
		}
	}
}

// authenticateAndSaveAPIKey authenticates with the bridge and saves the API key
func authenticateAndSaveAPIKey(cfg *config.Config, configPath string) {
	log.Printf("Setting up authentication with Hue bridge at %s", cfg.Hue.BridgeIP)
	log.Println("🔗 Press the link button on your Hue bridge now...")

	apiKey, err := hue.AuthenticateWithBridge(cfg.Hue.BridgeIP)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		log.Println("You can manually set the api_key in config.json or run again to retry authentication")
		return
	}

	log.Printf("\n✅ Authentication successful!")
	// Safe string truncation
	keyPreview := apiKey
	if len(apiKey) > 10 {
		keyPreview = apiKey[:10] + "..."
	}
	log.Printf("API key obtained: %s", keyPreview)

	// Update config with the new API key
	cfg.Hue.APIKey = apiKey

	// Save the updated configuration
	if err := config.SaveConfig(cfg, configPath); err != nil {
		log.Printf("Warning: Failed to save updated config: %v", err)
	} else {
		log.Println("Configuration saved with new API key")
	}
}
