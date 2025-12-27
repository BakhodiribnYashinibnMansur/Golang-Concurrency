package main

import "errors"

// Publish sends a message to all subscribers of a specific topic.
// This implements the broadcast pattern where one message is delivered to multiple subscribers.
//
// Go Concurrency Patterns used:
//   - RLock (Read Lock): Uses read lock because we're only reading the map, not modifying it
//     This allows multiple publishers to publish concurrently to different topics
//   - Channel send operation: Uses ch <- message to send message to each subscriber channel
//   - Non-blocking send: If channel is buffered and has capacity, send won't block
//   - Broadcast pattern: One message sent to multiple channels (fan-out pattern)
//
// Concurrency characteristics:
//   - Multiple publishers can publish to different topics concurrently (RLock allows this)
//   - If multiple publishers publish to the same topic, they'll serialize on the RLock
//   - Channel sends may block if subscriber channels are full (unbuffered or full buffer)
//
// Parameters:
//   - topic: string - the topic name to publish to
//   - message: string - the message content to broadcast
//
// Returns:
//   - error: returns error if topic doesn't exist
//
// Note: If a subscriber's channel is full, the send operation will block until space is available.
// This is a design choice - it ensures no messages are lost, but may slow down publishers.
func (p *Publisher) Publish(topic string, message string) error {
	p.RLock()         // Acquire read lock (allows concurrent reads, blocks writes)
	defer p.RUnlock() // Ensure lock is released

	// Get list of subscribers for this topic
	subscriber, ok := p.subscribers[topic]
	if !ok {
		return errors.New("topic not found")
	}

	// Broadcast message to all subscribers (fan-out pattern)
	// Each subscriber receives the message through their dedicated channel
	for _, ch := range subscriber {
		ch <- message // Send message to subscriber's channel
	}
	return nil
}
