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

	fmt.Println("Testing separated all/on and all/brightness handlers...")
	time.Sleep(2 * time.Second)

	// Test 1: All lights on
	fmt.Println("\n1. Turning all lights on...")
	msg := osc.NewMessage("/hue/all/on")
	msg.Append(int32(1)) // on = true
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 2: All lights brightness
	fmt.Println("\n2. Setting all lights brightness to 50%...")
	msg = osc.NewMessage("/hue/all/brightness")
	msg.Append(float32(0.5)) // 50% brightness
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 3: All lights on with transition
	fmt.Println("\n3. Turning all lights off with 2000ms transition...")
	msg = osc.NewMessage("/hue/all/on")
	msg.Append(int32(0))    // off = false
	msg.Append(int32(2000)) // 2 second transition
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 4: All lights brightness with transition
	fmt.Println("\n4. Setting all lights brightness to 80% with 1500ms transition...")
	msg = osc.NewMessage("/hue/all/brightness")
	msg.Append(float32(0.8)) // 80% brightness
	msg.Append(int32(1500))  // 1.5 second transition
	client.Send(msg)

	fmt.Println("\nTest completed! Check the logs:")
	fmt.Println("- All/on should show: 'All lights turned true/false'")
	fmt.Println("- All/brightness should show: 'All lights updated: [brightness=X%]'")
}
