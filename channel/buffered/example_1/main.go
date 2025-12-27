package main

import (
	"fmt"
	"time"
)

// worker simulates a worker that processes a task.
// It receives a value from the channel when done to signal completion.
//
// Parameters:
//   - id: worker identifier
//   - ch: buffered channel used for signaling completion
func worker(id int, ch chan int) {
	// Pretend we're doing some work
	fmt.Printf("Worker %d started\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d finished\n", id)
	// Receive from channel to signal completion
	// This removes one value from the buffer, allowing another worker to start
	<-ch // Signal we're done
}

// main demonstrates using buffered channels for worker pool pattern.
//
// Buffered Channel Characteristics:
//   - Created with make(chan Type, capacity) - capacity > 0
//   - Asynchronous when buffer has space: sender doesn't block if buffer isn't full
//   - Sender blocks only when buffer is full
//   - Receiver blocks only when buffer is empty
//   - Can store up to 'capacity' values without blocking
//
// Worker Pool Pattern:
//   - Limits concurrent workers to channel capacity
//   - Sending to channel before starting worker reserves a slot
//   - Worker receives from channel when done, freeing the slot
//   - This pattern controls resource usage and prevents too many concurrent operations
//
// Go Concurrency Pattern:
//   - Rate limiting: Buffer capacity limits concurrent operations
//   - Resource management: Channel acts as semaphore
//   - Backpressure: When buffer is full, new workers wait
//
// Flow:
//   1. Create buffered channel with capacity 3
//   2. Start 5 workers, but only 3 can run concurrently (buffer size)
//   3. First 3 workers start immediately (buffer has space)
//   4. Workers 4 and 5 wait until buffer has space
//   5. As workers finish, they free buffer slots, allowing waiting workers to start
func main() {
	// Create a buffered channel with capacity 3
	// This allows up to 3 values to be stored without blocking
	bufferedChannel := make(chan int, 3)

	// Start 5 workers
	// Only 3 can run concurrently due to buffer capacity
	for i := 1; i <= 5; i++ {
		// Start worker goroutine
		go worker(i, bufferedChannel)
		
		// Send to channel: reserves a slot in the buffer
		// First 3 sends complete immediately (buffer has space)
		// 4th and 5th sends block until buffer has space (when workers finish)
		bufferedChannel <- i // Fill the channel slot to represent a worker in-progress
	}

	// Wait for all workers to finish
	// This is a simple way to wait; in real-world scenarios, you might use sync.WaitGroup or similar
	// 7 seconds is enough for all 5 workers (each takes 1 second, but only 3 run concurrently)
	time.Sleep(7 * time.Second)
}
