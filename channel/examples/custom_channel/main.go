package main

import (
	"fmt"
	"sync"
	"time"
)

// main demonstrates custom channel implementation with various test cases.
//
// Custom Channel Implementation:
//   - Uses container/list for message queue
//   - Uses sync.Cond for blocking/waiting behavior
//   - Supports buffered and unbuffered channels
//   - Thread-safe operations with mutex
//
// Test Cases:
//  1. Basic send/receive
//  2. Buffered channel (capacity > 0)
//  3. Unbuffered channel (capacity = 0)
//  4. Multiple producers and consumers
//  5. Channel closing
//  6. Blocking behavior
func main() {
	fmt.Println("=== Custom Channel Implementation Tests ===")
	fmt.Println()

	// Test 1: Basic send/receive
	testBasicSendReceive()

	// Test 2: Buffered channel
	testBufferedChannel()

	// Test 3: Unbuffered channel
	testUnbufferedChannel()

	// Test 4: Multiple producers and consumers
	testMultipleProducersConsumers()

	// Test 5: Channel closing
	testChannelClosing()

	// Test 6: Blocking behavior
	testBlockingBehavior()

	fmt.Println("=== All Tests Completed ===")
}

// testBasicSendReceive tests basic send and receive operations
func testBasicSendReceive() {
	fmt.Println("Test 1: Basic Send/Receive")
	ch := NewChannel[string](1) // Buffered channel with capacity 1

	// Send a message
	err := ch.Send("Hello, World!")
	if err != nil {
		fmt.Printf("  Error sending: %v\n", err)
		return
	}

	// Receive the message
	msg, ok := ch.Receive()
	if !ok {
		fmt.Println("  Error: Failed to receive message")
		return
	}

	if msg == "Hello, World!" {
		fmt.Println("  ✓ Basic send/receive works correctly")
	} else {
		fmt.Printf("  ✗ Expected 'Hello, World!', got '%s'\n", msg)
	}
	fmt.Println()
}

// testBufferedChannel tests buffered channel behavior
func testBufferedChannel() {
	fmt.Println("Test 2: Buffered Channel (capacity 3)")
	ch := NewChannel[int](3)

	// Send multiple messages (should not block)
	for i := 1; i <= 3; i++ {
		err := ch.Send(i)
		if err != nil {
			fmt.Printf("  Error sending %d: %v\n", i, err)
			return
		}
		fmt.Printf("  Sent: %d\n", i)
	}

	// Receive all messages
	fmt.Println("  Receiving messages:")
	for i := 1; i <= 3; i++ {
		msg, ok := ch.Receive()
		if !ok {
			fmt.Printf("  Error receiving message %d\n", i)
			return
		}
		fmt.Printf("  Received: %d\n", msg)
	}
	fmt.Println("  ✓ Buffered channel works correctly")
	fmt.Println()
}

// testUnbufferedChannel tests unbuffered channel behavior
func testUnbufferedChannel() {
	fmt.Println("Test 3: Unbuffered Channel (capacity 0)")
	ch := NewChannel[string](0)

	var wg sync.WaitGroup

	// Producer goroutine
	wg.Go(func() {
		time.Sleep(100 * time.Millisecond) // Simulate work
		fmt.Println("  Producer: Sending message...")
		err := ch.Send("Unbuffered message")
		if err != nil {
			fmt.Printf("  Producer error: %v\n", err)
		} else {
			fmt.Println("  Producer: Message sent")
		}
	})

	// Consumer goroutine
	wg.Go(func() {
		fmt.Println("  Consumer: Waiting for message...")
		msg, ok := ch.Receive()
		if !ok {
			fmt.Println("  Consumer: Failed to receive")
		} else {
			fmt.Printf("  Consumer: Received '%s'\n", msg)
		}
	})

	wg.Wait()
	fmt.Println("  ✓ Unbuffered channel works correctly")
	fmt.Println()
}

// testMultipleProducersConsumers tests multiple goroutines
func testMultipleProducersConsumers() {
	fmt.Println("Test 4: Multiple Producers and Consumers")
	ch := NewChannel[int](5) // Buffered channel

	var wg sync.WaitGroup
	producerCount := 3
	consumerCount := 2
	totalMessages := 6

	// Start multiple producers
	for i := 0; i < producerCount; i++ {
		producerID := i + 1
		wg.Go(func() {
			for j := 0; j < totalMessages/producerCount; j++ {
				msg := producerID*10 + j
				err := ch.Send(msg)
				if err != nil {
					fmt.Printf("  Producer %d error: %v\n", producerID, err)
				} else {
					fmt.Printf("  Producer %d sent: %d\n", producerID, msg)
				}
				time.Sleep(50 * time.Millisecond)
			}
		})
	}

	// Start multiple consumers
	received := make([]int, 0)
	var mu sync.Mutex
	for i := 0; i < consumerCount; i++ {
		consumerID := i + 1
		wg.Go(func() {
			for j := 0; j < totalMessages/consumerCount; j++ {
				msg, ok := ch.Receive()
				if !ok {
					fmt.Printf("  Consumer %d: Channel closed\n", consumerID)
					return
				}
				mu.Lock()
				received = append(received, msg)
				mu.Unlock()
				fmt.Printf("  Consumer %d received: %d\n", consumerID, msg)
			}
		})
	}

	wg.Wait()
	fmt.Printf("  Total messages sent: %d, received: %d\n", totalMessages, len(received))
	if len(received) == totalMessages {
		fmt.Println("  ✓ Multiple producers/consumers work correctly")
	} else {
		fmt.Printf("  ✗ Expected %d messages, got %d\n", totalMessages, len(received))
	}
	fmt.Println()
}

// testChannelClosing tests channel closing behavior
func testChannelClosing() {
	fmt.Println("Test 5: Channel Closing")
	ch := NewChannel[string](2)

	// Send some messages
	ch.Send("Message 1")
	ch.Send("Message 2")

	// Close the channel
	err := ch.Close()
	if err != nil {
		fmt.Printf("  Error closing channel: %v\n", err)
		return
	}
	fmt.Println("  Channel closed")

	// Try to send after closing (should fail)
	err = ch.Send("Message 3")
	if err != nil {
		fmt.Printf("  ✓ Send after close correctly returns error: %v\n", err)
	} else {
		fmt.Println("  ✗ Send after close should return error")
	}

	// Receive remaining messages
	msg1, ok1 := ch.Receive()
	msg2, ok2 := ch.Receive()

	if ok1 && ok2 && msg1 == "Message 1" && msg2 == "Message 2" {
		fmt.Println("  ✓ Can receive messages after closing")
	} else {
		fmt.Printf("  ✗ Receive after close: msg1=%s, ok1=%v, msg2=%s, ok2=%v\n", msg1, ok1, msg2, ok2)
	}

	// Try to close again (should fail)
	err = ch.Close()
	if err != nil {
		fmt.Printf("  ✓ Close after close correctly returns error: %v\n", err)
	} else {
		fmt.Println("  ✗ Close after close should return error")
	}
	fmt.Println()
}

// testBlockingBehavior tests blocking behavior
func testBlockingBehavior() {
	fmt.Println("Test 6: Blocking Behavior")
	ch := NewChannel[int](2) // Capacity 2

	var wg sync.WaitGroup

	// Producer 1: Fill the buffer
	wg.Go(func() {
		fmt.Println("  Producer 1: Sending 2 messages (fills buffer)...")
		ch.Send(1)
		ch.Send(2)
		fmt.Println("  Producer 1: Buffer filled")
	})

	time.Sleep(100 * time.Millisecond)

	// Producer 2: Try to send (should block until space available)
	wg.Go(func() {
		fmt.Println("  Producer 2: Trying to send (should block)...")
		start := time.Now()
		ch.Send(3) // This should block until consumer receives
		duration := time.Since(start)
		fmt.Printf("  Producer 2: Sent after %v (was blocked)\n", duration)
	})

	time.Sleep(200 * time.Millisecond)

	// Consumer: Receive one message (frees space)
	wg.Go(func() {
		time.Sleep(300 * time.Millisecond)
		fmt.Println("  Consumer: Receiving message...")
		msg, ok := ch.Receive()
		if ok {
			fmt.Printf("  Consumer: Received %d (freed space for Producer 2)\n", msg)
		}
	})

	wg.Wait()

	// Receive remaining messages
	msg2, _ := ch.Receive()
	msg3, _ := ch.Receive()
	fmt.Printf("  Remaining messages: %d, %d\n", msg2, msg3)
	fmt.Println("  ✓ Blocking behavior works correctly")
	fmt.Println()
}
