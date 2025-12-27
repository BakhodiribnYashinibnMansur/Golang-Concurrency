package main

import (
	"fmt"
	"time"
)

// main demonstrates channel closing and checking if channel is open.
//
// Channel Closing:
//   - close(ch) signals that no more values will be sent
//   - Closing a channel sends a zero value to all waiting receivers
//   - Receivers can check if channel is closed using two-value receive: value, ok := <-ch
//   - Sending to a closed channel causes panic
//   - Closing an already closed channel causes panic
//
// Two-Value Receive:
//   - value, ok := <-ch
//   - ok is true if value was received, false if channel is closed and empty
//   - When channel is closed, ok becomes false and value is zero value
//
// Go Concurrency Pattern:
//   - Graceful shutdown: Close channel to signal completion
//   - Channel state checking: Verify channel is open before processing
//   - Resource cleanup: Use defer to ensure channel is closed
//
// Flow:
//   1. Create unbuffered channel
//   2. Start goroutine that sends message after 2 seconds
//   3. Defer channel closing (executes when function exits)
//   4. Receive with two-value form to check if channel is open
//   5. If channel is closed, print message and return
//   6. Otherwise, print received message
func main() {
	// Create an unbuffered channel
	messageChannel := make(chan string)
	
	// Start goroutine that sends message after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		messageChannel <- "Hello World"
	}()
	
	// Defer channel closing: ensures channel is closed when function exits
	// This is important for cleanup, though in this example the channel
	// will be closed before the goroutine sends (which would cause panic)
	// In real code, close after all sends are complete
	defer close(messageChannel)
	
	// Two-value receive: checks if channel is open
	// message: the received value (or zero value if channel closed)
	// open: true if value received, false if channel is closed
	message, open := <-messageChannel
	if !open {
		fmt.Println("Channel closed")
		return
	}
	fmt.Println(message)
}
