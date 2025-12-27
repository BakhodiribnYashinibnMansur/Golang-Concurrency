package main

import (
	"fmt"
	"time"
)

// main demonstrates basic unbuffered channel communication.
//
// Unbuffered Channel Characteristics:
//   - Created with make(chan Type) - no capacity specified
//   - Synchronous communication: sender blocks until receiver is ready
//   - Receiver blocks until sender sends a value
//   - Direct handoff: value is transferred directly from sender to receiver
//   - No buffer: cannot store values, only passes them directly
//
// Go Concurrency Pattern:
//   - Goroutine communication: Uses channel for message passing between goroutines
//   - Blocking send/receive: Both operations block until the other side is ready
//   - Synchronization: Channel acts as synchronization primitive
//
// Flow:
//  1. Create unbuffered channel (no buffer capacity)
//  2. Start goroutine that sleeps 2 seconds, then sends message
//  3. Main goroutine blocks on receive until message arrives
//  4. When goroutine sends, value is directly transferred to main goroutine
//  5. Both goroutines continue execution
func main() {
	// Create an unbuffered channel (capacity = 0)
	// This channel can only pass values when both sender and receiver are ready
	messageChannel := make(chan string)

	// Start a goroutine that will send a message after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		// Send operation blocks until receiver is ready
		// Since main goroutine is already waiting, this send completes immediately
		messageChannel <- "Hello World"
	}()

	// Receive operation blocks until sender sends a value
	// Main goroutine waits here for 2 seconds until goroutine sends the message
	message := <-messageChannel
	fmt.Println(message)
}
