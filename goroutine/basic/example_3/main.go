package main

import (
	"fmt"
	"sync"
	"time"
)

// worker simulates a worker that processes a task.
//
// Parameters:
//   - id: worker identifier
//   - wg: WaitGroup pointer for synchronization
func worker(id int, wg *sync.WaitGroup) {
	// Decrement counter when goroutine completes
	// defer ensures Done() is called even if function panics
	defer wg.Done()

	fmt.Printf("Worker %d: Starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d: Finished\n", id)
}

// main demonstrates using sync.WaitGroup for goroutine synchronization.
//
// sync.WaitGroup Characteristics:
//   - Counter-based synchronization primitive
//   - Add(n): Increments counter by n
//   - Done(): Decrements counter by 1
//   - Wait(): Blocks until counter reaches 0
//   - Safe for concurrent use
//
// Go Concurrency Pattern:
//   - Synchronization: WaitGroup ensures all goroutines complete
//   - No sleep needed: Wait() blocks until all workers finish
//   - Proper cleanup: Main waits for all workers before exiting
//
// Flow:
//  1. Create WaitGroup
//  2. For each worker: Add(1) to increment counter
//  3. Start worker goroutine (passes WaitGroup pointer)
//  4. Worker calls Done() when finished (decrements counter)
//  5. Main calls Wait() to block until counter is 0
//  6. All workers complete, main continues
func main() {
	var wg sync.WaitGroup

	fmt.Println("Main: Starting workers with WaitGroup")

	// Start 5 workers
	for i := 1; i <= 5; i++ {
		wg.Add(1) // Increment counter before starting goroutine
		go worker(i, &wg)
	}

	fmt.Println("Main: Waiting for workers to finish...")

	// Block until all workers call Done()
	// This is better than time.Sleep() because it waits exactly as long as needed
	wg.Wait()

	fmt.Println("Main: All workers completed")
}
