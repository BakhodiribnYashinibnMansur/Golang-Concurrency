package main

import (
	"fmt"
	"time"
)

// main demonstrates using select statement with multiple unbuffered channels.
//
// Select Statement:
//   - Allows waiting on multiple channel operations simultaneously
//   - Executes the first case that becomes ready
//   - If multiple cases are ready, one is chosen randomly
//   - Blocks until at least one case is ready
//
// Go Concurrency Pattern:
//   - Multiplexing: Select allows handling multiple channels concurrently
//   - Non-blocking when possible: Can use default case for non-blocking operations
//   - Channel selection: Choose which channel to receive from based on availability
//
// Flow:
//  1. Create two unbuffered channels
//  2. Start two goroutines that send values after different delays (1s and 2s)
//  3. Use select to receive from whichever channel is ready first
//  4. First select receives from ch1 (after 1 second)
//  5. Second select receives from ch2 (after 2 seconds)
func main() {
	// Create two unbuffered channels
	ch1 := make(chan int)
	ch2 := make(chan int)

	// Goroutine 1: sends value 1 after 1 second
	go func() {
		time.Sleep(time.Second * 1)
		ch1 <- 1 // Send blocks until receiver is ready
	}()

	// Goroutine 2: sends value 2 after 2 seconds
	go func() {
		time.Sleep(time.Second * 2)
		ch2 <- 2 // Send blocks until receiver is ready
	}()

	// Loop twice to receive from both channels
	// for range 2 is Go 1.22+ syntax (iterates 2 times)
	for range 2 {
		// Select statement: waits for whichever channel is ready first
		// First iteration: receives from ch1 (ready after 1 second)
		// Second iteration: receives from ch2 (ready after 2 seconds)
		select {
		case msg1 := <-ch1:
			fmt.Println("Received from ch1:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received from ch2:", msg2)
		}
		// Note: If both channels are ready, Go randomly selects one
	}
}
