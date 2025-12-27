package main

import (
	"fmt"
	"sync"
)

// safeCounter demonstrates proper synchronization with mutex.
//
// Parameters:
//   - count: shared counter variable (pointer)
//   - mu: mutex for protecting shared variable
//   - wg: WaitGroup for synchronization
func safeCounter(count *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 1000; i++ {
		// Lock before accessing shared variable
		mu.Lock()
		*count++ // Critical section: only one goroutine can execute this at a time
		mu.Unlock()
	}
}

// main demonstrates fixing race condition with sync.Mutex.
//
// sync.Mutex Characteristics:
//   - Mutual exclusion lock
//   - Lock(): Acquires lock (blocks if already locked)
//   - Unlock(): Releases lock
//   - Only one goroutine can hold lock at a time
//   - Protects critical sections
//
// Go Concurrency Pattern:
//   - Synchronization: Mutex ensures only one goroutine accesses shared data
//   - Critical section: Code between Lock() and Unlock()
//   - Thread-safe: Prevents race conditions
//
// Flow:
//  1. Create mutex and shared counter
//  2. Start 10 goroutines, all incrementing counter
//  3. Each goroutine locks mutex before increment
//  4. Only one goroutine can increment at a time
//  5. After increment, goroutine unlocks mutex
//  6. Final count is correct: 10,000
//
// Comparison with example_4:
//   - example_4: No synchronization → race condition → incorrect result
//   - example_5: Mutex synchronization → no race condition → correct result
func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	count := 0

	fmt.Println("Main: Starting goroutines with mutex protection")
	fmt.Println("This code is thread-safe!")
	fmt.Println()

	// Start 10 goroutines, all incrementing the same counter
	// But now with mutex protection
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go safeCounter(&count, &mu, &wg)
	}

	wg.Wait()

	fmt.Printf("Expected count: 10000\n")
	fmt.Printf("Actual count:   %d\n", count)

	if count == 10000 {
		fmt.Println("✅ Success! Count is correct.")
		fmt.Println("   Mutex protected shared variable from race condition.")
	} else {
		fmt.Println("❌ Unexpected error! Count should be 10000.")
	}
}
