package main

import (
	"fmt"
	"sync"
	"time"
)

// main demonstrates a complete Pub/Sub system implementation using Go concurrency patterns.
//
// Go Concurrency Patterns demonstrated:
//  1. Goroutines: Multiple concurrent publishers and subscribers
//  2. Channels: Message passing between publishers and subscribers
//  3. sync.WaitGroup: Coordinating multiple goroutines (using wg.Go() from Go 1.21+)
//  4. RWMutex: Thread-safe access to shared state (in Publisher struct)
//  5. Channel range: Receiving messages until channel is closed
//  6. Channel closing: Graceful shutdown pattern
//  7. Fan-out pattern: One publisher broadcasts to multiple subscribers
//
// Architecture:
//   - Publisher-Subscriber pattern: Decouples message producers from consumers
//   - Topic-based routing: Messages are organized by topics
//   - Broadcast semantics: Each message is delivered to all subscribers of a topic
//   - Concurrent processing: Publishers and subscribers run in separate goroutines
//
// Flow:
//  1. Create topics and initialize Publisher
//  2. Start publisher goroutines that publish messages to topics
//  3. Start subscriber goroutines that subscribe to topics and receive messages
//  4. Wait for all goroutines to complete using WaitGroup
//  5. Gracefully close all topics, which closes subscriber channels
func main() {
	// Define all topics and their configuration in one place
	// This centralizes configuration and makes it easy to add/modify topics
	topicConfig := map[string]struct {
		messages        []string      // Messages to publish for this topic
		subscriberCount int           // Number of subscribers for this topic
		delay           time.Duration // Delay between message publications
	}{
		"news": {
			messages: []string{
				"Breaking news: Major announcement!",
				"Breaking news: Important update!",
				"Breaking news: Latest development!",
			},
			subscriberCount: 3,
			delay:           500 * time.Millisecond,
		},
		"sports": {
			messages: []string{
				"Game score: 3-1",
				"Match update: Half-time",
				"Final score: 4-2",
			},
			subscriberCount: 2,
			delay:           600 * time.Millisecond,
		},
		"tech": {
			messages: []string{
				"New technology released!",
				"Tech update: Version 2.0",
				"Tech news: Innovation breakthrough!",
			},
			subscriberCount: 1,
			delay:           550 * time.Millisecond,
		},
	}

	// Create a new Publisher instance
	// Publisher uses RWMutex internally for thread-safe operations
	pub := NewPublisher()

	// Create all topics before starting publishers/subscribers
	// Topics must exist before subscribers can subscribe or publishers can publish
	topics := make([]string, 0, len(topicConfig))
	for topic := range topicConfig {
		topics = append(topics, topic)
		pub.CreateTopic(topic)
		fmt.Printf("Topic '%s' created\n", topic)
	}

	// WaitGroup coordinates multiple goroutines
	// We use separate WaitGroups for publishers and subscribers
	// because subscribers wait for channels to close, while publishers finish after sending messages
	var publisherWg sync.WaitGroup
	var subscriberWg sync.WaitGroup

	// Start multiple publisher goroutines
	// Each publisher runs in its own goroutine and publishes messages to a topic
	fmt.Println("\n=== Starting Publishers ===")

	for topic, config := range topicConfig {
		// Capture loop variables to avoid closure issues
		topicName := topic
		messages := config.messages
		delay := config.delay

		// wg.Go() starts a goroutine and automatically increments WaitGroup counter
		// When goroutine completes, it should call wg.Done() (but wg.Go handles this)
		publisherWg.Go(func() {
			// Publisher loop: publish each message with delay
			for i, msg := range messages {
				time.Sleep(delay) // Simulate work/delay between publications

				// Publish message to topic (broadcasts to all subscribers)
				if err := pub.Publish(topicName, msg); err != nil {
					fmt.Printf("Publisher error publishing to [%s]: %v\n", topicName, err)
				} else {
					fmt.Printf("Publisher → [%s]: %s\n", topicName, msg)
				}

				// Extra delay after last message
				if i == len(messages)-1 {
					time.Sleep(200 * time.Millisecond)
				}
			}
		})
	}

	// Start multiple subscriber goroutines
	// Each subscriber runs in its own goroutine and receives messages from a topic
	fmt.Println("\n=== Starting Subscribers ===")

	for topic, config := range topicConfig {
		topicName := topic
		subscriberCount := config.subscriberCount

		// Create multiple subscribers for each topic (demonstrates broadcast pattern)
		for i := 1; i <= subscriberCount; i++ {
			subID := i

			// wg.Go() starts subscriber goroutine
			subscriberWg.Go(func() {
				// Subscribe to topic and get receive-only channel
				ch, err := pub.Subscribe(topicName)
				if err != nil {
					fmt.Printf("Subscriber %d error subscribing to [%s]: %v\n", subID, topicName, err)
					return
				}
				fmt.Printf("Subscriber %d subscribed to [%s]\n", subID, topicName)

				// Range over channel: receives messages until channel is closed
				// This is the idiomatic Go pattern for receiving from channels
				// When channel is closed, loop exits automatically
				for msg := range ch {
					// Process received message
					fmt.Printf("  Subscriber %d ← [%s]: %s\n", subID, topicName, msg)
				}

				// This line executes when channel is closed (graceful shutdown)
				fmt.Printf("Subscriber %d unsubscribed from [%s]\n", subID, topicName)
			})
		}
	}

	// Wait a bit for subscribers to register before publishers start sending
	// This ensures subscribers are ready to receive messages
	time.Sleep(100 * time.Millisecond)

	// Wait for all publisher goroutines to finish
	// Publishers finish after sending all their messages
	fmt.Println("\n=== Waiting for publishers to finish ===")
	publisherWg.Wait()

	// Wait a bit more for any remaining messages to be processed
	time.Sleep(500 * time.Millisecond)

	// Gracefully close all topics
	// Closing a topic closes all subscriber channels, which causes
	// "for msg := range ch" loops to exit (graceful shutdown pattern)
	fmt.Println("\n=== Closing topics ===")
	for _, topic := range topics {
		if err := pub.CloseTopic(topic); err != nil {
			fmt.Printf("Error closing topic '%s': %v\n", topic, err)
		} else {
			fmt.Printf("Topic '%s' closed\n", topic)
		}
	}

	// Wait for all subscriber goroutines to finish
	// Subscribers finish when their channels are closed (after topics are closed)
	fmt.Println("\n=== Waiting for subscribers to finish ===")
	subscriberWg.Wait()

	fmt.Println("\n=== Program completed ===")
}
