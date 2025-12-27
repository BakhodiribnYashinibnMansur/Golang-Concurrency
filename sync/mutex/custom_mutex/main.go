package main

import (
	"fmt"
	"sync"
	"time"
)

// main demonstrates custom mutex (monitor) implementation with various test cases.
//
// Custom Mutex Implementation:
//   - Uses channels to manage state (monitor pattern)
//   - Thread-safe data access without sync.Mutex
//   - Serializes access via a dedicated goroutine
//
// Test Cases:
//  1. Basic Operations
//  2. Concurrent Access
//  3. Resource Clean-up
func main() {
	fmt.Println("=== Custom Mutex Implementation Tests ===")
	fmt.Println()

	// Test 1: Basic Operations
	testBasicOperations()

	// Test 2: Concurrent Access
	testConcurrentAccess()

	// Test 3: Resource Clean-up
	testCleanup()

	// Test 4: Heavy Load (1000 values)
	testManyValues()

	// Test 5: Concurrent (Race Condition) String Access
	testConcurrentValueAccess()

	fmt.Println("=== All Tests Completed ===")
}

// testBasicOperations tests simple Get and Send operations
func testBasicOperations() {
	fmt.Println("Test 1: Basic Operations")
	m := NewMutex[int]()

	// Send a value
	fmt.Println("  Send(10)...")
	m.Send(10)

	// Get the value
	val := m.Get()
	fmt.Printf("  Get() returned: %d\n", val)

	if val == 10 {
		fmt.Println("  ✓ Value matches expected (10)")
	} else {
		fmt.Printf("  ✗ Expected 10, got %d\n", val)
	}

	// Update value
	fmt.Println("  Send(55)...")
	m.Send(55)

	val = m.Get()
	if val == 55 {
		fmt.Println("  ✓ Value updated correctly to 55")
	} else {
		fmt.Printf("  ✗ Expected 55, got %d\n", val)
	}

	m.Close()
	fmt.Println()
}

// testConcurrentAccess tests safety under high concurrency
func testConcurrentAccess() {
	fmt.Println("Test 2: Concurrent Access (Stress Test)")
	m := NewMutex[int]()
	var wg sync.WaitGroup

	writerCount := 10
	readerCount := 10
	iterations := 100

	fmt.Printf("  Starting %d writers and %d readers (%d iterations each)...\n", writerCount, readerCount, iterations)

	// Writers
	for i := 0; i < writerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				m.Send(j)
				// Small sleep to allow context switching
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	// Readers
	for i := 0; i < readerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = m.Get()
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	// Wait for all to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("  ✓ Operations completed without deadlock")
	case <-time.After(2 * time.Second):
		fmt.Println("  ✗ Timeout: Operations took too long (possible deadlock)")
	}

	m.Close()
	fmt.Println()
}

// testCleanup tests closing the mutex
func testCleanup() {
	fmt.Println("Test 3: Resource Clean-up")
	m := NewMutex[int]()

	m.Send(1)
	fmt.Println("  Mutex running...")

	m.Close()
	fmt.Println("  Mutex closed")

	// Note: In this implementation, calling Send/Get after Close
	// would block forever because the monitor loop has exited
	// and no one is reading from the channels.
	// We just verify here that Close() doesn't panic.

	fmt.Println("  ✓ Close() successful")
	fmt.Println()
}

// testManyValues tests the mutex with 1000 distinct sequential updates
func testManyValues() {
	fmt.Println("Test 4: Heavy Load (1000 distinct values)")

	// Use NewMutexWithValue to start with a specific state
	initialVal := -1
	m := NewMutexWithValue(initialVal)

	val := m.Get()
	if val == initialVal {
		fmt.Printf("  ✓ Initial value correct: %d\n", val)
	} else {
		fmt.Printf("  ✗ Expected initial value %d, got %d\n", initialVal, val)
	}

	fmt.Println("  Running 1000 updates...")
	start := time.Now()

	success := true
	for i := 0; i < 1000; i++ {
		m.Send(i)
		got := m.Get()
		if got != i {
			fmt.Printf("  ✗ Failed at %d: expected %d, got %d\n", i, i, got)
			success = false
			break
		}
	}

	if success {
		fmt.Printf("  ✓ Successfully processed 1000 values in %v\n", time.Since(start))
	}

	m.Close()
	fmt.Println()
}

// testConcurrentStringAccess tests race conditions with string type
func testConcurrentValueAccess() {
	fmt.Println("Test 5: Concurrent String Access (Race Condition Simulation)")
	m := NewMutex[string]()
	var wg sync.WaitGroup

	writers := 100
	readers := 100
	ops := 100

	fmt.Printf("  Starting %d writers and %d readers (%d ops each) with strings...\n", writers, readers, ops)

	// Writers
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ops; j++ {
				val := fmt.Sprintf("writer-%d-iter-%d", id, j)
				m.Send(val)
				// Small randomization to mix up schedule
				if j%10 == 0 {
					time.Sleep(time.Microsecond)
				}
			}
		}(i)
	}

	// Readers
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ops; j++ {
				_ = m.Get() // Just consume
				if j%10 == 0 {
					time.Sleep(time.Microsecond)
				}
			}
		}(i)
	}

	// Wrap wait in a channel to detect deadlocks/timeouts
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("  ✓ String operations completed without deadlock")
	case <-time.After(3 * time.Second):
		fmt.Println("  ✗ Timeout: String operations took too long")
	}

	m.Close()
	fmt.Println()
}
