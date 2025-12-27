package main

import (
	"fmt"
)

// main demonstrates sending values to buffered channel, closing it, and receiving all values.
//
// Buffered Channel Behavior:
//   - Can send multiple values without blocking (up to capacity)
//   - Values are stored in buffer until received
//   - Sender doesn't block if buffer has space
//
// Range Over Channel:
//   - for value := range ch iterates over channel values
//   - Automatically receives values until channel is closed
//   - Loop exits when channel is closed and empty
//   - Only works with closed channels or channels that will be closed
//
// Go Concurrency Pattern:
//   - Batch processing: Send multiple values, then close and receive all
//   - Producer-consumer: Producer sends to buffer, consumer receives
//   - Channel iteration: Idiomatic way to receive all values from channel
//
// Flow:
//   1. Create buffered channel with capacity 3
//   2. Send 3 values to channel (fills buffer, no blocking)
//   3. Close channel to signal no more values
//   4. Range over channel receives all 3 values
//   5. Loop exits when channel is closed and empty
func main() {
	// Create a buffered channel with capacity 3
	bufferedChannel := make(chan int, 3)

	// Send 3 values to the channel
	// All sends complete immediately because buffer has space
	// No blocking occurs since we're not exceeding capacity
	for i := 1; i <= 3; i++ {
		bufferedChannel <- i // Fill the channel slot to represent a worker in-progress
	}
	
	// Close the channel to signal that no more values will be sent
	// Important: Close after all sends are complete
	// Closing allows receivers to know when all values have been sent
	close(bufferedChannel)
	
	// Range over the channel to receive values until it is closed
	// This is the idiomatic Go pattern for receiving all values
	// Loop automatically exits when channel is closed and empty
	for i := range bufferedChannel {
		fmt.Println("Received value:", i)
	}
	// After loop exits, all values have been received
}
