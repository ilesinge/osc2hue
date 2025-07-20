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

	fmt.Println("Testing all/on handler delegation...")
	time.Sleep(2 * time.Second)

	// Test 1: All lights on
	fmt.Println("\n1. Turning all lights on...")
	msg := osc.NewMessage("/hue/all/on")
	msg.Append(int32(1)) // on = true
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 2: All lights off with transition
	fmt.Println("\n2. Turning all lights off with 2000ms transition...")
	msg = osc.NewMessage("/hue/all/on")
	msg.Append(int32(0))    // off = false
	msg.Append(int32(2000)) // 2 second transition
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 3: All lights on again (no transition)
	fmt.Println("\n3. Turning all lights on again...")
	msg = osc.NewMessage("/hue/all/on")
	msg.Append(int32(1)) // on = true
	client.Send(msg)

	fmt.Println("\nTest completed! Check the logs:")
	fmt.Println("- Should show individual light messages: 'Light {id} turned true/false'")
	fmt.Println("- Should show summary message: 'All lights turned true/false'")
	fmt.Println("- Each light should handle its own transition duration")
}
