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

	fmt.Println("Testing color handler (now using set handler internally)...")
	time.Sleep(2 * time.Second)

	// Test 1: Individual light color without transition
	fmt.Println("\n1. Setting light 1 color to warm white (x=0.4, y=0.4)...")
	msg := osc.NewMessage("/hue/light/1/color")
	msg.Append(float32(0.4)) // x coordinate
	msg.Append(float32(0.4)) // y coordinate
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 2: Individual light color with transition
	fmt.Println("\n2. Setting light 1 color to cool blue (x=0.15, y=0.06) with 2000ms transition...")
	msg = osc.NewMessage("/hue/light/1/color")
	msg.Append(float32(0.15)) // x coordinate
	msg.Append(float32(0.06)) // y coordinate
	msg.Append(int32(2000))   // 2 second transition
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 3: Another light with different color
	fmt.Println("\n3. Setting light 2 color to red (x=0.7, y=0.3)...")
	msg = osc.NewMessage("/hue/light/2/color")
	msg.Append(float32(0.7)) // x coordinate
	msg.Append(float32(0.3)) // y coordinate
	client.Send(msg)
	time.Sleep(3 * time.Second)

	// Test 4: Color with slow transition
	fmt.Println("\n4. Setting light 3 color to green (x=0.3, y=0.6) with 3000ms transition...")
	msg = osc.NewMessage("/hue/light/3/color")
	msg.Append(float32(0.3)) // x coordinate
	msg.Append(float32(0.6)) // y coordinate
	msg.Append(int32(3000))  // 3 second transition
	client.Send(msg)

	fmt.Println("\nTest completed! The color handler should now delegate to set handler.")
	fmt.Println("Check the logs - you should see 'color=x:X,y:Y' in the output, indicating")
	fmt.Println("that the set handler is being used internally.")
}
