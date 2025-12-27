package main

import "errors"

// CloseTopic closes all subscriber channels for a topic and removes it from the Publisher.
// This implements graceful shutdown pattern for a topic.
//
// Go Concurrency Patterns used:
//   - Channel closing: close(ch) signals to receivers that no more values will be sent
//   - Range over channel pattern: When a channel is closed, "for msg := range ch" loops exit
//   - Graceful shutdown: Allows subscribers to finish processing current messages before exiting
//
// Channel closing semantics:
//   - Closing a channel sends a zero value to all waiting receivers
//   - Receivers using "for msg := range ch" will exit the loop when channel is closed
//   - Sending to a closed channel causes panic, but closing is safe to call multiple times
//
// Parameters:
//   - topic: string - the topic name to close
//
// Returns:
//   - error: returns error if topic doesn't exist
//
// Usage: Call this when you want to stop a topic and notify all subscribers to stop listening.
func (p *Publisher) CloseTopic(topic string) error {
	p.Lock()         // Acquire exclusive write lock (modifying map)
	defer p.Unlock() // Ensure lock is released

	if _, ok := p.subscribers[topic]; !ok {
		return errors.New("topic not found")
	}

	// Close all subscriber channels for this topic
	// This causes all "for msg := range ch" loops in subscribers to exit
	for _, ch := range p.subscribers[topic] {
		close(ch) // Signal no more messages will be sent
	}

	// Remove topic from map
	delete(p.subscribers, topic)
	return nil
}

// CloseSubscriber removes a specific subscriber from a topic by closing their channel
// and removing it from the subscriber list.
//
// Go Concurrency Patterns used:
//   - Channel closing: Closes the specific subscriber's channel
//   - Slice manipulation: Removes channel from slice using append with slice slicing
//   - Selective shutdown: Allows removing individual subscribers without closing entire topic
//
// Parameters:
//   - topic: string - the topic name
//   - subscriberChannel: <-chan string - the subscriber's channel to remove
//
// Returns:
//   - error: returns error if topic or subscriber not found
//
// Note: This method closes the channel, which will cause the subscriber's range loop to exit.
func (p *Publisher) CloseSubscriber(topic string, subscriberChannel <-chan string) error {
	p.Lock()         // Acquire exclusive write lock
	defer p.Unlock() // Ensure lock is released

	if _, ok := p.subscribers[topic]; !ok {
		return errors.New("topic not found")
	}

	// Find and remove the subscriber's channel from the list
	for i, subscriber := range p.subscribers[topic] {
		// Compare channels (receive-only channel can be compared with bidirectional channel)
		if subscriber == subscriberChannel {
			// Close the bidirectional channel stored in map (not the receive-only parameter)
			// This signals the subscriber that no more messages will be sent
			close(subscriber)

			// Remove channel from slice using slice slicing
			p.subscribers[topic] = append(p.subscribers[topic][:i], p.subscribers[topic][i+1:]...)
			return nil
		}
	}
	return errors.New("subscriber not found")
}
