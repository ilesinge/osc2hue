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

	fmt.Println("Testing brightness handlers (now using set handlers internally)...")
	time.Sleep(2 * time.Second)

	// Test 1: Individual light brightness
	fmt.Println("\n1. Setting light 1 brightness to 50% (using brightness handler)...")
	msg := osc.NewMessage("/hue/1/brightness")
	msg.Append(float32(0.5)) // 50% brightness
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 2: Individual light brightness with transition
	fmt.Println("\n2. Setting light 1 brightness to 80% with 2000ms transition...")
	msg = osc.NewMessage("/hue/1/brightness")
	msg.Append(float32(0.8)) // 80% brightness
	msg.Append(int32(2000))  // 2 second transition
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 3: All lights brightness
	fmt.Println("\n3. Setting all lights brightness to 30% (using brightness handler)...")
	msg = osc.NewMessage("/hue/all/brightness")
	msg.Append(float32(0.3)) // 30% brightness
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 4: All lights brightness with transition
	fmt.Println("\n4. Setting all lights brightness to 70% with 1500ms transition...")
	msg = osc.NewMessage("/hue/all/brightness")
	msg.Append(float32(0.7)) // 70% brightness
	msg.Append(int32(1500))  // 1.5 second transition
	client.Send(msg)

	fmt.Println("\nTest completed! The brightness handlers should now delegate to set handlers.")
	fmt.Println("Check the logs - you should see 'brightness=X%' in the output, indicating")
	fmt.Println("that the set handlers are being used internally.")
}
