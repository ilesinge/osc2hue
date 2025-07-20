//go:build example
// +build example

package main

import (
	"fmt"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

func main() {
	// Create OSC client
	client := osc.NewClient("localhost", 8080)

	fmt.Println("🎨 Testing /hue/all/color handler...")
	fmt.Println("Watch your lights change colors!")
	time.Sleep(2 * time.Second)

	// Test 1: Warm white (instant)
	fmt.Println("\n1. Setting all lights to warm white (instant)")
	msg := osc.NewMessage("/hue/all/color")
	msg.Append(float32(0.4)) // x coordinate
	msg.Append(float32(0.4)) // y coordinate
	client.Send(msg)
	time.Sleep(2 * time.Second)

	// Test 2: Cool blue with transition
	fmt.Println("\n2. Setting all lights to cool blue (1 second transition)")
	msg = osc.NewMessage("/hue/all/color")
	msg.Append(float32(0.15)) // x coordinate
	msg.Append(float32(0.06)) // y coordinate
	msg.Append(int32(1000))   // 1000ms transition
	client.Send(msg)
	time.Sleep(2 * time.Second)

	// Test 3: Warm red
	fmt.Println("\n3. Setting all lights to warm red (500ms transition)")
	msg = osc.NewMessage("/hue/all/color")
	msg.Append(float32(0.6)) // x coordinate
	msg.Append(float32(0.3)) // y coordinate
	msg.Append(int32(500))   // 500ms transition
	client.Send(msg)
	time.Sleep(1500 * time.Millisecond)

	// Test 4: Green
	fmt.Println("\n4. Setting all lights to green (2 second transition)")
	msg = osc.NewMessage("/hue/all/color")
	msg.Append(float32(0.3)) // x coordinate
	msg.Append(float32(0.6)) // y coordinate
	msg.Append(int32(2000))  // 2000ms transition
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 5: Purple/Magenta
	fmt.Println("\n5. Setting all lights to purple (instant)")
	msg = osc.NewMessage("/hue/all/color")
	msg.Append(float32(0.25)) // x coordinate
	msg.Append(float32(0.1))  // y coordinate
	client.Send(msg)
	time.Sleep(2 * time.Second)

	// Test 6: Yellow/Orange
	fmt.Println("\n6. Setting all lights to yellow/orange (1.5 second transition)")
	msg = osc.NewMessage("/hue/all/color")
	msg.Append(float32(0.5))  // x coordinate
	msg.Append(float32(0.45)) // y coordinate
	msg.Append(int32(1500))   // 1500ms transition
	client.Send(msg)
	time.Sleep(2500 * time.Millisecond)

	// Test 7: Back to neutral warm white
	fmt.Println("\n7. Returning all lights to neutral warm white (3 second fade)")
	msg = osc.NewMessage("/hue/all/color")
	msg.Append(float32(0.35)) // x coordinate
	msg.Append(float32(0.35)) // y coordinate
	msg.Append(int32(3000))   // 3000ms transition
	client.Send(msg)
	time.Sleep(4 * time.Second)

	fmt.Println("\n✅ All color tests completed!")
	fmt.Println("Did you see all the color changes on your lights?")
}
