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

	fmt.Println("Testing /set command with null values...")
	fmt.Println("Note: These commands use -1 for null/skip parameters")

	// Wait for user to see the message
	time.Sleep(2 * time.Second)

	// Test 1: Set only color (x=0.3, y=0.4), skip brightness and duration
	fmt.Println("\n1. Setting only color (x=0.3, y=0.4), skipping brightness...")
	msg := osc.NewMessage("/hue/all/set")
	msg.Append(float32(0.3)) // x coordinate
	msg.Append(float32(0.4)) // y coordinate
	msg.Append(int32(-1))    // brightness = null (skip)
	msg.Append(int32(-1))    // duration = null (skip)
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 2: Set only brightness (50%), skip color and duration
	fmt.Println("\n2. Setting only brightness (50%), skipping color...")
	msg = osc.NewMessage("/hue/all/set")
	msg.Append(int32(-1))    // x = null (skip)
	msg.Append(int32(-1))    // y = null (skip)
	msg.Append(float32(0.5)) // brightness = 50%
	msg.Append(int32(-1))    // duration = null (skip)
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 3: Set color + brightness with transition duration
	fmt.Println("\n3. Setting color (x=0.6, y=0.3) + brightness (80%) with 2000ms transition...")
	msg = osc.NewMessage("/hue/all/set")
	msg.Append(float32(0.6)) // x coordinate
	msg.Append(float32(0.3)) // y coordinate
	msg.Append(float32(0.8)) // brightness = 80%
	msg.Append(int32(2000))  // duration = 2000ms
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 4: Set only transition duration (fade current state over 3 seconds)
	fmt.Println("\n4. Setting only transition duration (3000ms), keeping current color/brightness...")
	msg = osc.NewMessage("/hue/all/set")
	msg.Append(int32(-1))   // x = null (skip)
	msg.Append(int32(-1))   // y = null (skip)
	msg.Append(int32(-1))   // brightness = null (skip)
	msg.Append(int32(3000)) // duration = 3000ms
	client.Send(msg)
	time.Sleep(4 * time.Second)

	// Test 5: Individual light test - set only brightness for light 1
	fmt.Println("\n5. Setting only brightness (30%) for light 1...")
	msg = osc.NewMessage("/hue/1/set")
	msg.Append(int32(-1))    // x = null (skip)
	msg.Append(int32(-1))    // y = null (skip)
	msg.Append(float32(0.3)) // brightness = 30%
	client.Send(msg)

	fmt.Println("\nTest completed! Check the osc2hue logs to see the null parameter handling.")
}
