package main

import (
	"fmt"
	"sync"
)

// counter demonstrates the race condition problem.
//
// Race Condition:
//   - Multiple goroutines access shared variable without synchronization
//   - Read-modify-write operations are not atomic
//   - Final result is unpredictable and incorrect
//
// Parameters:
//   - count: shared counter variable (pointer)
//   - wg: WaitGroup for synchronization
func counter(count *int, wg *sync.WaitGroup) {
	defer wg.Done()

	// This is NOT thread-safe!
	// Multiple goroutines read and write 'count' simultaneously
	// Each increment involves: read current value, add 1, write new value
	// These steps can interleave, causing lost updates
	for i := 0; i < 1000; i++ {
		*count++ // Race condition here!
	}
}

// main demonstrates race condition with shared variable.
//
// Race Condition Problem:
//   - Multiple goroutines modify shared variable concurrently
//   - No synchronization mechanism (mutex, channel, atomic)
//   - Operations are not atomic
//   - Results are non-deterministic
//
// Go Concurrency Pattern (Anti-pattern):
//   - This is an example of INCORRECT concurrent code
//   - Shows why synchronization is necessary
//   - Run with 'go run -race main.go' to detect race condition
//
// Expected vs Actual:
//   - Expected: 10 goroutines × 1000 increments = 10,000
//   - Actual: Usually less than 10,000 due to lost updates
//   - Each run produces different result
//
// Flow:
//  1. Start 10 goroutines, all incrementing shared counter
//  2. Goroutines run concurrently, causing race condition
//  3. Some increments are lost due to concurrent access
//  4. Final count is incorrect and non-deterministic
//
// Solution:
//   - Use sync.Mutex to protect shared variable
//   - Use channels for communication
//   - Use atomic operations (sync/atomic package)
func main() {
	var wg sync.WaitGroup
	count := 0

	fmt.Println("Main: Starting goroutines with race condition")
	fmt.Println("WARNING: This code has a race condition!")
	fmt.Println("Run with: go run -race main.go to detect it")
	fmt.Println()

	// Start 10 goroutines, all incrementing the same counter
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go counter(&count, &wg)
	}

	wg.Wait()

	fmt.Printf("Expected count: 10000\n")
	fmt.Printf("Actual count:   %d\n", count)

	if count != 10000 {
		fmt.Println("❌ Race condition detected! Count is incorrect.")
		fmt.Println("   Multiple goroutines modified shared variable without synchronization.")
	} else {
		fmt.Println("⚠️  Count is correct this time, but race condition still exists!")
		fmt.Println("   Run multiple times or use -race flag to detect it.")
	}
}
