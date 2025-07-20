package main

import (
	"fmt"
	"log"
	"strings"

	"osc2hue/internal/osc"

	gosc "github.com/hypebeast/go-osc/osc"
	"github.com/openhue/openhue-go"
)

// addAllHandlers adds all OSC handlers (individual lights and global commands)
func addAllHandlers(oscServer *osc.Server, home *openhue.Home, lights []openhue.LightGet) {
	// Add individual light handlers
	addLightHandlers(oscServer, home, lights)

	// Add global handlers
	addGlobalHandlers(oscServer, home, lights)
}

// addGlobalHandlers adds OSC handlers for global "all lights" commands
func addGlobalHandlers(oscServer *osc.Server, home *openhue.Home, lights []openhue.LightGet) {
	oscServer.AddHandler("/hue/all/on", func(msg *gosc.Message) {
		handleAllOn(msg, home, lights)
	})

	oscServer.AddHandler("/hue/all/brightness", func(msg *gosc.Message) {
		handleAllBrightness(msg, home, lights)
	})

	oscServer.AddHandler("/hue/all/color", func(msg *gosc.Message) {
		handleAllColor(msg, home, lights)
	})

	oscServer.AddHandler("/hue/all/set", func(msg *gosc.Message) {
		handleAllSet(msg, home, lights)
	})
}

// addLightHandlers adds OSC handlers for all discovered lights
func addLightHandlers(oscServer *osc.Server, home *openhue.Home, lights []openhue.LightGet) {
	if home == nil {
		return
	}

	for _, light := range lights {
		// Convert light ID to string for the closure
		lightID := *light.Id

		// Create handlers with proper closure capture
		oscServer.AddHandler(fmt.Sprintf("/hue/light/%s/on", lightID), func(msg *gosc.Message) {
			handleLightOn(msg, home, lightID)
		})

		oscServer.AddHandler(fmt.Sprintf("/hue/light/%s/brightness", lightID), func(msg *gosc.Message) {
			handleLightBrightness(msg, home, lightID)
		})

		oscServer.AddHandler(fmt.Sprintf("/hue/light/%s/color", lightID), func(msg *gosc.Message) {
			handleLightColor(msg, home, lightID)
		})

		// Combined color+brightness handler
		oscServer.AddHandler(fmt.Sprintf("/hue/light/%s/set", lightID), func(msg *gosc.Message) {
			handleLightSet(msg, home, lightID)
		})
	}

	// Also add numeric handlers for convenience (1, 2, 3, etc.)
	for i, light := range lights {
		lightID := *light.Id
		numericID := i + 1

		oscServer.AddHandler(fmt.Sprintf("/hue/light/%d/on", numericID), func(msg *gosc.Message) {
			handleLightOn(msg, home, lightID)
		})

		oscServer.AddHandler(fmt.Sprintf("/hue/light/%d/brightness", numericID), func(msg *gosc.Message) {
			handleLightBrightness(msg, home, lightID)
		})

		oscServer.AddHandler(fmt.Sprintf("/hue/light/%d/color", numericID), func(msg *gosc.Message) {
			handleLightColor(msg, home, lightID)
		})

		// Combined color+brightness handler for numeric IDs
		oscServer.AddHandler(fmt.Sprintf("/hue/light/%d/set", numericID), func(msg *gosc.Message) {
			handleLightSet(msg, home, lightID)
		})
	}
}

func handleLightOn(msg *gosc.Message, home *openhue.Home, lightID string) {
	if home == nil {
		log.Printf("Hue bridge not connected")
		return
	}

	if len(msg.Arguments) < 1 {
		log.Printf("No arguments provided for light on/off")
		return
	}

	var on bool
	switch v := msg.Arguments[0].(type) {
	case int32:
		on = v > 0
	case float32:
		on = v > 0
	case bool:
		on = v
	default:
		log.Printf("Invalid argument type for light on/off: %T", v)
		return
	}

	// Create light state update
	state := openhue.LightPut{
		On: &openhue.On{On: &on},
	}

	// Check if transition duration is provided as second argument
	if len(msg.Arguments) >= 2 {
		var transitionMs int
		switch v := msg.Arguments[1].(type) {
		case int32:
			transitionMs = int(v)
		case float32:
			transitionMs = int(v)
		default:
			log.Printf("Invalid transition duration type: %T", v)
		}
		if transitionMs > 0 {
			state.Dynamics = &openhue.LightDynamics{Duration: &transitionMs}
		}
	}

	if err := home.UpdateLight(lightID, state); err != nil {
		log.Printf("Error setting light state: %v", err)
	} else {
		log.Printf("Light %s turned %v", lightID, on)
	}
}

func handleLightBrightness(msg *gosc.Message, home *openhue.Home, lightID string) {
	if home == nil {
		log.Printf("Hue bridge not connected")
		return
	}

	if len(msg.Arguments) < 1 {
		log.Printf("No arguments provided for brightness")
		return
	}

	// Create a new message for the set handler with null color values
	setMsg := gosc.NewMessage("/hue/light/set")
	setMsg.Append(int32(-1))        // x = null (skip color)
	setMsg.Append(int32(-1))        // y = null (skip color)
	setMsg.Append(msg.Arguments[0]) // brightness value from original message

	// Add transition duration if provided
	if len(msg.Arguments) >= 2 {
		setMsg.Append(msg.Arguments[1]) // duration from original message
	}

	// Delegate to the set handler
	handleLightSet(setMsg, home, lightID)
}

func handleLightColor(msg *gosc.Message, home *openhue.Home, lightID string) {
	if home == nil {
		log.Printf("Hue bridge not connected")
		return
	}

	if len(msg.Arguments) < 2 {
		log.Printf("Not enough arguments for color (need X and Y coordinates)")
		return
	}

	// Create a new message for the set handler with null brightness value
	setMsg := gosc.NewMessage("/hue/light/set")
	setMsg.Append(msg.Arguments[0]) // x coordinate from original message
	setMsg.Append(msg.Arguments[1]) // y coordinate from original message
	setMsg.Append(int32(-1))        // brightness = null (skip)

	// Add transition duration if provided
	if len(msg.Arguments) >= 3 {
		setMsg.Append(msg.Arguments[2]) // duration from original message
	}

	// Delegate to the set handler
	handleLightSet(setMsg, home, lightID)
}

func handleLightSet(msg *gosc.Message, home *openhue.Home, lightID string) {
	if home == nil {
		log.Printf("Hue bridge not connected")
		return
	}

	// Allow flexible number of arguments, but require at least 1
	if len(msg.Arguments) < 1 {
		log.Printf("Set command requires at least 1 argument. Use -1 for null values.")
		log.Printf("Usage: /hue/light/{id}/set {x|-1} [y|-1] [brightness|-1] [duration_ms|-1]")
		return
	}

	// Create light state update
	state := openhue.LightPut{}

	var hasColor, hasBrightness bool
	var x, y float64
	var brightness float64
	var logParts []string

	// Parse X coordinate (argument 0)
	if len(msg.Arguments) >= 1 {
		switch v := msg.Arguments[0].(type) {
		case int32:
			if v != -1 {
				x = float64(v)
				hasColor = true
			}
		case float32:
			if v != -1.0 {
				x = float64(v)
				hasColor = true
			}
		default:
			log.Printf("Invalid X coordinate type: %T", v)
			return
		}
	}

	// Parse Y coordinate (argument 1)
	if len(msg.Arguments) >= 2 && hasColor {
		switch v := msg.Arguments[1].(type) {
		case int32:
			if v != -1 {
				y = float64(v)
			} else {
				hasColor = false // If Y is null, disable color
			}
		case float32:
			if v != -1.0 {
				y = float64(v)
			} else {
				hasColor = false // If Y is null, disable color
			}
		default:
			log.Printf("Invalid Y coordinate type: %T", v)
			return
		}
	} else if hasColor && len(msg.Arguments) < 2 {
		log.Printf("Color requires both X and Y coordinates")
		return
	}

	// Parse brightness (argument 2)
	if len(msg.Arguments) >= 3 {
		switch v := msg.Arguments[2].(type) {
		case int32:
			if v != -1 {
				if v <= 254 {
					brightness = float64(v) / 254.0 // 0-254 range
				} else {
					brightness = float64(v) / 100.0 // assume percentage
				}
				hasBrightness = true
			}
		case float32:
			if v != -1.0 {
				if v <= 1.0 {
					brightness = float64(v) // 0.0-1.0 range
				} else if v <= 100.0 {
					brightness = float64(v) / 100.0 // percentage
				} else {
					brightness = float64(v) / 254.0 // 0-254 range
				}
				hasBrightness = true
			}
		default:
			log.Printf("Invalid brightness type: %T", v)
			return
		}
	}

	// Apply color if specified
	if hasColor {
		// Validate and clamp values
		if x < 0 {
			x = 0
		} else if x > 1 {
			x = 1
		}
		if y < 0 {
			y = 0
		} else if y > 1 {
			y = 1
		}

		// Convert coordinates to float32 for the API
		xf := float32(x)
		yf := float32(y)

		state.Color = &openhue.Color{
			Xy: &openhue.GamutPosition{
				X: &xf,
				Y: &yf,
			},
		}
		logParts = append(logParts, fmt.Sprintf("color=x:%.3f,y:%.3f", x, y))
	}

	// Apply brightness if specified
	if hasBrightness {
		// Validate and clamp brightness
		if brightness < 0 {
			brightness = 0
		} else if brightness > 1 {
			brightness = 1
		}

		// Convert brightness to percentage for OpenHue
		brightnessPercent := float32(brightness * 100)
		on := brightness > 0

		state.On = &openhue.On{On: &on}
		state.Dimming = &openhue.Dimming{Brightness: &brightnessPercent}
		logParts = append(logParts, fmt.Sprintf("brightness=%.1f%%", brightnessPercent))
	}

	// Check if transition duration is provided as fourth argument
	if len(msg.Arguments) >= 4 {
		var transitionMs int
		switch v := msg.Arguments[3].(type) {
		case int32:
			if v != -1 {
				transitionMs = int(v)
			}
		case float32:
			if v != -1.0 {
				transitionMs = int(v)
			}
		default:
			log.Printf("Invalid transition duration type: %T", v)
		}
		if transitionMs > 0 {
			state.Dynamics = &openhue.LightDynamics{Duration: &transitionMs}
			logParts = append(logParts, fmt.Sprintf("duration=%dms", transitionMs))
		}
	}

	// Check if we have anything to update
	if !hasColor && !hasBrightness && state.Dynamics == nil {
		log.Printf("No valid parameters provided for light %s", lightID)
		return
	}

	if err := home.UpdateLight(lightID, state); err != nil {
		log.Printf("Error updating light %s: %v", lightID, err)
	} else {
		if len(logParts) > 0 {
			log.Printf("Light %s updated: %s", lightID, fmt.Sprintf("[%s]", strings.Join(logParts, ", ")))
		} else {
			log.Printf("Light %s updated (no changes applied)", lightID)
		}
	}
}

func handleAllOn(msg *gosc.Message, home *openhue.Home, lights []openhue.LightGet) {
	if home == nil {
		log.Printf("Hue bridge not connected")
		return
	}

	if len(msg.Arguments) < 1 {
		log.Printf("No arguments provided for all lights on/off")
		return
	}

	// Extract the on/off value for logging
	var on bool
	switch v := msg.Arguments[0].(type) {
	case int32:
		on = v > 0
	case float32:
		on = v > 0
	case bool:
		on = v
	default:
		log.Printf("Invalid argument type for all lights on/off: %T", v)
		return
	}

	// Apply to all lights using the single light handler
	for _, light := range lights {
		handleLightOn(msg, home, *light.Id)
	}
	log.Printf("All lights turned %v", on)
}

func handleAllBrightness(msg *gosc.Message, home *openhue.Home, lights []openhue.LightGet) {
	if home == nil {
		log.Printf("Hue bridge not connected")
		return
	}

	if len(msg.Arguments) < 1 {
		log.Printf("No arguments provided for all lights brightness")
		return
	}

	// Create a new message for the set handler with null color values
	setMsg := gosc.NewMessage("/hue/all/set")
	setMsg.Append(int32(-1))        // x = null (skip color)
	setMsg.Append(int32(-1))        // y = null (skip color)
	setMsg.Append(msg.Arguments[0]) // brightness value from original message

	// Add transition duration if provided
	if len(msg.Arguments) >= 2 {
		setMsg.Append(msg.Arguments[1]) // duration from original message
	}

	// Delegate to the set handler
	handleAllSet(setMsg, home, lights)
}

func handleAllColor(msg *gosc.Message, home *openhue.Home, lights []openhue.LightGet) {
	if home == nil {
		log.Printf("Hue bridge not connected")
		return
	}

	if len(msg.Arguments) < 2 {
		log.Printf("Not enough arguments for all lights color (need X and Y coordinates)")
		return
	}

	// Create a new message for the set handler with null brightness value
	setMsg := gosc.NewMessage("/hue/all/set")
	setMsg.Append(msg.Arguments[0]) // x coordinate from original message
	setMsg.Append(msg.Arguments[1]) // y coordinate from original message
	setMsg.Append(int32(-1))        // brightness = null (skip)

	// Add transition duration if provided
	if len(msg.Arguments) >= 3 {
		setMsg.Append(msg.Arguments[2]) // duration from original message
	}

	// Delegate to the set handler
	handleAllSet(setMsg, home, lights)
}

func handleAllSet(msg *gosc.Message, home *openhue.Home, lights []openhue.LightGet) {
	if home == nil {
		log.Printf("Hue bridge not connected")
		return
	}

	// Allow flexible number of arguments, but require at least 1
	if len(msg.Arguments) < 1 {
		log.Printf("Set command requires at least 1 argument. Use -1 for null values.")
		log.Printf("Usage: /hue/all/set {x|-1} [y|-1] [brightness|-1] [duration_ms|-1]")
		return
	}

	// Apply to all lights using the single light handler
	for _, light := range lights {
		handleLightSet(msg, home, *light.Id)
	}
	log.Printf("All lights updated")
}
