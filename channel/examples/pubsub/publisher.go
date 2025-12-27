package main

import "sync"

// Publisher implements the Publisher-Subscriber (Pub/Sub) pattern using Go channels.
// This is a concurrent-safe message broker that allows multiple publishers to send
// messages to topics and multiple subscribers to receive them.
//
// Go Concurrency Patterns used:
//   - Channel-based communication: Uses channels for message passing between goroutines
//   - RWMutex: Read-Write mutex for efficient concurrent access (multiple readers, single writer)
//   - Thread-safe map: Protects shared state (subscribers map) from race conditions
//
// Architecture:
//   - Each topic maintains a slice of channels (one per subscriber)
//   - When a message is published, it's sent to all subscriber channels (broadcast pattern)
//   - Subscribers receive messages through their dedicated channel
type Publisher struct {
	sync.RWMutex                          // Protects subscribers map from concurrent access
	subscribers  map[string][]chan string // Topic -> list of subscriber channels
}

// NewPublisher creates and returns a new Publisher instance.
// Initializes the subscribers map to store topic-channel mappings.
//
// Returns: *Publisher - pointer to the newly created Publisher
func NewPublisher() *Publisher {
	return &Publisher{
		subscribers: make(map[string][]chan string),
	}
}
