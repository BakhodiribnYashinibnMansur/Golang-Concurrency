package main

import (
	"fmt"
	"time"
)

// worker simulates a worker that processes tasks.
//
// Parameters:
//   - id: worker identifier for logging
//   - tasks: number of tasks this worker should process
func worker(id int, tasks int) {
	for i := 1; i <= tasks; i++ {
		fmt.Printf("Worker %d: Processing task %d\n", id, i)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Printf("Worker %d: Completed all tasks\n", id)
}

// main demonstrates multiple goroutines running concurrently.
//
// Goroutine Characteristics:
//   - Multiple goroutines can run simultaneously
//   - Each goroutine has independent execution flow
//   - Goroutines are scheduled by Go runtime (M:N scheduling)
//   - Output order is non-deterministic (depends on scheduler)
//
// Go Concurrency Pattern:
//   - Parallel execution: Multiple workers process tasks simultaneously
//   - Independent execution: Each worker runs independently
//   - Concurrent output: Messages from different workers interleave
//
// Flow:
//  1. Main goroutine starts 3 worker goroutines
//  2. Each worker processes 3 tasks concurrently
//  3. Workers run in parallel, output interleaves
//  4. Main waits for all workers to complete
//  5. Program exits when main goroutine finishes
func main() {
	fmt.Println("Main: Starting multiple workers")

	// Start 3 workers, each processing 3 tasks
	// All workers run concurrently
	for i := 1; i <= 3; i++ {
		go worker(i, 3)
	}

	// Wait for all workers to complete
	// 3 workers × 3 tasks × 500ms = ~4.5 seconds total
	// But they run concurrently, so actual time is ~1.5 seconds
	time.Sleep(2 * time.Second)

	fmt.Println("Main: All workers finished")
}
