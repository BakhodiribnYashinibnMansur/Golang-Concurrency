package main

import "errors"

// Subscribe allows a subscriber to register for messages from a specific topic.
// Returns a receive-only channel (<-chan string) that the subscriber can use to receive messages.
//
// Go Concurrency Patterns used:
//   - Channel creation: Creates a buffered channel (capacity 1) for the subscriber
//   - Receive-only channel: Returns <-chan string to prevent subscribers from sending
//   - Channel-based communication: Messages flow through channels between goroutines
//   - Lock for map modification: Uses exclusive lock to safely append to subscribers slice
//
// Channel Pattern:
//   - Buffered channel (capacity 1): Allows one message to be buffered, preventing blocking
//   - Range over channel: Subscribers use "for msg := range ch" to receive messages
//   - Channel closing: When topic is closed, all subscriber channels are closed, causing
//     range loops to exit gracefully
//
// Parameters:
//   - topic: string - the topic name to subscribe to
//
// Returns:
//   - <-chan string: receive-only channel for receiving messages
//   - error: returns error if topic doesn't exist
//
// Usage example:
//
//		ch, err := pub.Subscribe("news")
//		if err != nil { ... }
//		for msg := range ch {
//	 	Process message
//		}
func (p *Publisher) Subscribe(topic string) (<-chan string, error) {
	p.Lock()         // Acquire exclusive write lock (modifying subscribers map)
	defer p.Unlock() // Ensure lock is released

	// Create buffered channel with capacity 1 for this subscriber
	// Buffered channel prevents blocking if subscriber is slow to read
	channel := make(chan string, 1)

	// Check if topic exists
	if _, ok := p.subscribers[topic]; !ok {
		return nil, errors.New("topic not found")
	}

	// Add subscriber's channel to the topic's subscriber list
	p.subscribers[topic] = append(p.subscribers[topic], channel)
	return channel, nil
}
