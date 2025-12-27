package main

import (
	"fmt"
)

// main demonstrates producer-consumer pattern with buffered channel.
//
// Producer-Consumer Pattern:
//   - Producer goroutine sends values to channel
//   - Consumer goroutine receives values from channel
//   - Buffered channel decouples producer and consumer
//   - Producer can send multiple values without waiting for consumer
//
// Buffered Channel Advantages:
//   - Non-blocking sends: Producer doesn't block if buffer has space
//   - Smooth operation: Buffer absorbs temporary speed differences
//   - Better performance: Reduces goroutine blocking
//
// Go Concurrency Pattern:
//   - Goroutine communication: Producer and consumer run concurrently
//   - Channel closing: Producer closes channel when done
//   - Range over channel: Consumer receives all values until channel closed
//   - Synchronization: Channel coordinates producer and consumer
//
// Flow:
//   1. Create buffered channel with capacity 3
//   2. Start producer goroutine that sends 5 values
//   3. Producer sends first 3 values immediately (buffer has space)
//   4. Producer blocks on 4th send until consumer receives (buffer full)
//   5. Consumer receives values as they become available
//   6. Producer closes channel after sending all values
//   7. Consumer's range loop exits when channel is closed and empty
func main() {
	// Create a buffered channel with capacity 3
	bufferedChannel := make(chan int, 3)

	// Producer goroutine: sends values to channel
	go func() {
		// Send 5 values (more than buffer capacity)
		for i := 1; i <= 5; i++ {
			// First 3 sends complete immediately (buffer has space)
			// 4th and 5th sends block until consumer receives (buffer full)
			bufferedChannel <- i
		}
		// Close channel after all values are sent
		// This signals consumer that no more values will come
		close(bufferedChannel)
	}()

	// Consumer: receives all values from channel
	// Range loop automatically receives values until channel is closed
	for i := range bufferedChannel {
		fmt.Println("Received value:", i)
	}
	// Loop exits when channel is closed and all values received
}
