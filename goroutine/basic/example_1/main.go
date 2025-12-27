package main

import (
	"fmt"
	"time"
)

// main demonstrates basic goroutine creation and execution.
//
// Goroutine Characteristics:
//   - Lightweight threads managed by Go runtime
//   - Created with 'go' keyword before function call
//   - Runs concurrently with other goroutines
//   - Has its own call stack (starts small, grows as needed)
//   - Scheduled by Go runtime, not OS
//   - Very cheap to create (thousands can run simultaneously)
//
// Go Concurrency Pattern:
//   - Concurrent execution: Multiple goroutines run simultaneously
//   - Non-blocking: Main goroutine continues while worker runs
//   - Lightweight: Goroutines are much lighter than OS threads
//
// Flow:
//  1. Main goroutine starts
//  2. New goroutine is spawned with 'go' keyword
//  3. Both goroutines run concurrently
//  4. Worker goroutine prints message after 1 second
//  5. Main goroutine waits 2 seconds to ensure worker completes
//  6. Program exits when main goroutine finishes
func main() {
	// Start a goroutine
	// The 'go' keyword creates a new goroutine that runs concurrently
	go worker()

	// Main goroutine continues executing
	fmt.Println("Main: Started worker goroutine")

	// Wait for worker to finish
	// Without this sleep, main would exit before worker completes
	// In real code, use sync.WaitGroup or channels for synchronization
	time.Sleep(2 * time.Second)
	fmt.Println("Main: Exiting")
}

// worker simulates some work in a goroutine
func worker() {
	fmt.Println("Worker: Starting")
	time.Sleep(1 * time.Second)
	fmt.Println("Worker: Finished")
}
