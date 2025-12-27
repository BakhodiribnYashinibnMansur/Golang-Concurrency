package main

import (
	"sync"
	"testing"
	"time"
)

// TestNewPublisher tests the creation of a new Publisher instance
func TestNewPublisher(t *testing.T) {
	pub := NewPublisher()
	if pub == nil {
		t.Fatal("NewPublisher() returned nil")
	}
	if pub.subscribers == nil {
		t.Fatal("NewPublisher() subscribers map is nil")
	}
	if len(pub.subscribers) != 0 {
		t.Errorf("Expected empty subscribers map, got %d topics", len(pub.subscribers))
	}
}

// TestCreateTopic tests topic creation
func TestCreateTopic(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"

	pub.CreateTopic(topic)

	pub.RLock()
	defer pub.RUnlock()

	if _, ok := pub.subscribers[topic]; !ok {
		t.Fatal("Topic was not created")
	}
	if len(pub.subscribers[topic]) != 0 {
		t.Errorf("Expected empty subscriber list, got %d subscribers", len(pub.subscribers[topic]))
	}
}

// TestSubscribe tests subscribing to a topic
func TestSubscribe(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	ch, err := pub.Subscribe(topic)
	if err != nil {
		t.Fatalf("Subscribe() returned error: %v", err)
	}
	if ch == nil {
		t.Fatal("Subscribe() returned nil channel")
	}

	pub.RLock()
	if len(pub.subscribers[topic]) != 1 {
		t.Errorf("Expected 1 subscriber, got %d", len(pub.subscribers[topic]))
	}
	pub.RUnlock()
}

// TestSubscribeNonExistentTopic tests subscribing to a non-existent topic
func TestSubscribeNonExistentTopic(t *testing.T) {
	pub := NewPublisher()

	ch, err := pub.Subscribe("non-existent")
	if err == nil {
		t.Fatal("Expected error when subscribing to non-existent topic")
	}
	if ch != nil {
		t.Fatal("Expected nil channel when subscription fails")
	}
}

// TestPublish tests publishing messages to subscribers
func TestPublish(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	// Subscribe to topic
	ch, err := pub.Subscribe(topic)
	if err != nil {
		t.Fatalf("Subscribe() returned error: %v", err)
	}

	// Publish a message
	message := "test message"
	err = pub.Publish(topic, message)
	if err != nil {
		t.Fatalf("Publish() returned error: %v", err)
	}

	// Receive the message
	select {
	case msg := <-ch:
		if msg != message {
			t.Errorf("Expected message '%s', got '%s'", message, msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

// TestPublishMultipleSubscribers tests broadcasting to multiple subscribers
func TestPublishMultipleSubscribers(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	// Create multiple subscribers
	subscribers := make([]<-chan string, 3)
	for i := 0; i < 3; i++ {
		ch, err := pub.Subscribe(topic)
		if err != nil {
			t.Fatalf("Subscribe() returned error: %v", err)
		}
		subscribers[i] = ch
	}

	// Publish a message
	message := "broadcast message"
	err := pub.Publish(topic, message)
	if err != nil {
		t.Fatalf("Publish() returned error: %v", err)
	}

	// All subscribers should receive the message
	for i, ch := range subscribers {
		select {
		case msg := <-ch:
			if msg != message {
				t.Errorf("Subscriber %d: Expected message '%s', got '%s'", i, message, msg)
			}
		case <-time.After(1 * time.Second):
			t.Fatalf("Subscriber %d: Timeout waiting for message", i)
		}
	}
}

// TestPublishNonExistentTopic tests publishing to a non-existent topic
func TestPublishNonExistentTopic(t *testing.T) {
	pub := NewPublisher()

	err := pub.Publish("non-existent", "message")
	if err == nil {
		t.Fatal("Expected error when publishing to non-existent topic")
	}
}

// TestCloseTopic tests closing a topic
func TestCloseTopic(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	// Subscribe to topic
	ch, err := pub.Subscribe(topic)
	if err != nil {
		t.Fatalf("Subscribe() returned error: %v", err)
	}

	// Close the topic
	err = pub.CloseTopic(topic)
	if err != nil {
		t.Fatalf("CloseTopic() returned error: %v", err)
	}

	// Channel should be closed (receive should return zero value and ok=false)
	select {
	case msg, ok := <-ch:
		if ok {
			t.Errorf("Expected channel to be closed, but received message: %s", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for channel close")
	}

	// Topic should be removed
	pub.RLock()
	if _, ok := pub.subscribers[topic]; ok {
		t.Error("Topic should be removed after CloseTopic()")
	}
	pub.RUnlock()
}

// TestCloseTopicNonExistent tests closing a non-existent topic
func TestCloseTopicNonExistent(t *testing.T) {
	pub := NewPublisher()

	err := pub.CloseTopic("non-existent")
	if err == nil {
		t.Fatal("Expected error when closing non-existent topic")
	}
}

// TestCloseSubscriber tests closing a specific subscriber
func TestCloseSubscriber(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	// Subscribe to topic
	ch, err := pub.Subscribe(topic)
	if err != nil {
		t.Fatalf("Subscribe() returned error: %v", err)
	}

	// Close the subscriber
	err = pub.CloseSubscriber(topic, ch)
	if err != nil {
		t.Fatalf("CloseSubscriber() returned error: %v", err)
	}

	// Channel should be closed
	select {
	case msg, ok := <-ch:
		if ok {
			t.Errorf("Expected channel to be closed, but received message: %s", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for channel close")
	}

	// Subscriber should be removed from topic
	pub.RLock()
	if len(pub.subscribers[topic]) != 0 {
		t.Errorf("Expected 0 subscribers, got %d", len(pub.subscribers[topic]))
	}
	pub.RUnlock()
}

// TestCloseSubscriberNonExistentTopic tests closing subscriber from non-existent topic
func TestCloseSubscriberNonExistentTopic(t *testing.T) {
	pub := NewPublisher()
	ch := make(<-chan string)

	err := pub.CloseSubscriber("non-existent", ch)
	if err == nil {
		t.Fatal("Expected error when closing subscriber from non-existent topic")
	}
}

// TestCloseSubscriberNonExistentSubscriber tests closing non-existent subscriber
func TestCloseSubscriberNonExistentSubscriber(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	ch := make(<-chan string)

	err := pub.CloseSubscriber(topic, ch)
	if err == nil {
		t.Fatal("Expected error when closing non-existent subscriber")
	}
}

// TestConcurrentPublish tests concurrent publishing to the same topic
func TestConcurrentPublish(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	// Subscribe to topic
	ch, err := pub.Subscribe(topic)
	if err != nil {
		t.Fatalf("Subscribe() returned error: %v", err)
	}

	// Publish multiple messages concurrently
	var wg sync.WaitGroup
	numMessages := 10
	wg.Add(numMessages)

	for i := 0; i < numMessages; i++ {
		go func(msgNum int) {
			defer wg.Done()
			message := "concurrent message"
			if err := pub.Publish(topic, message); err != nil {
				t.Errorf("Publish() returned error: %v", err)
			}
		}(i)
	}

	wg.Wait()

	// Receive all messages
	received := make(map[string]int)
	timeout := time.After(2 * time.Second)
	for i := 0; i < numMessages; i++ {
		select {
		case msg := <-ch:
			received[msg]++
		case <-timeout:
			t.Fatalf("Timeout waiting for message %d", i)
		}
	}

	if received["concurrent message"] != numMessages {
		t.Errorf("Expected %d messages, got %d", numMessages, received["concurrent message"])
	}
}

// TestConcurrentSubscribe tests concurrent subscriptions
func TestConcurrentSubscribe(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	// Subscribe concurrently
	var wg sync.WaitGroup
	numSubscribers := 10
	channels := make([]<-chan string, numSubscribers)
	wg.Add(numSubscribers)

	for i := 0; i < numSubscribers; i++ {
		go func(idx int) {
			defer wg.Done()
			ch, err := pub.Subscribe(topic)
			if err != nil {
				t.Errorf("Subscribe() returned error: %v", err)
				return
			}
			channels[idx] = ch
		}(i)
	}

	wg.Wait()

	// Verify all subscribers were added
	pub.RLock()
	if len(pub.subscribers[topic]) != numSubscribers {
		t.Errorf("Expected %d subscribers, got %d", numSubscribers, len(pub.subscribers[topic]))
	}
	pub.RUnlock()

	// Publish a message and verify all subscribers receive it
	message := "broadcast to all"
	err := pub.Publish(topic, message)
	if err != nil {
		t.Fatalf("Publish() returned error: %v", err)
	}

	for i, ch := range channels {
		select {
		case msg := <-ch:
			if msg != message {
				t.Errorf("Subscriber %d: Expected message '%s', got '%s'", i, message, msg)
			}
		case <-time.After(1 * time.Second):
			t.Fatalf("Subscriber %d: Timeout waiting for message", i)
		}
	}
}

// TestMultipleTopics tests operations with multiple topics
func TestMultipleTopics(t *testing.T) {
	pub := NewPublisher()
	topics := []string{"topic1", "topic2", "topic3"}

	// Create topics
	for _, topic := range topics {
		pub.CreateTopic(topic)
	}

	// Subscribe to each topic
	channels := make(map[string]<-chan string)
	for _, topic := range topics {
		ch, err := pub.Subscribe(topic)
		if err != nil {
			t.Fatalf("Subscribe() to %s returned error: %v", topic, err)
		}
		channels[topic] = ch
	}

	// Publish to each topic
	for _, topic := range topics {
		message := "message for " + topic
		err := pub.Publish(topic, message)
		if err != nil {
			t.Fatalf("Publish() to %s returned error: %v", topic, err)
		}
	}

	// Verify each subscriber received the correct message
	for _, topic := range topics {
		expectedMsg := "message for " + topic
		select {
		case msg := <-channels[topic]:
			if msg != expectedMsg {
				t.Errorf("Topic %s: Expected message '%s', got '%s'", topic, expectedMsg, msg)
			}
		case <-time.After(1 * time.Second):
			t.Fatalf("Topic %s: Timeout waiting for message", topic)
		}
	}
}

// TestPublishAfterClose tests that publishing after closing a topic returns an error
func TestPublishAfterClose(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	// Close the topic
	err := pub.CloseTopic(topic)
	if err != nil {
		t.Fatalf("CloseTopic() returned error: %v", err)
	}

	// Try to publish after closing
	err = pub.Publish(topic, "message")
	if err == nil {
		t.Fatal("Expected error when publishing to closed topic")
	}
}

// TestSubscribeAfterClose tests that subscribing after closing a topic returns an error
func TestSubscribeAfterClose(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	// Close the topic
	err := pub.CloseTopic(topic)
	if err != nil {
		t.Fatalf("CloseTopic() returned error: %v", err)
	}

	// Try to subscribe after closing
	ch, err := pub.Subscribe(topic)
	if err == nil {
		t.Fatal("Expected error when subscribing to closed topic")
	}
	if ch != nil {
		t.Fatal("Expected nil channel when subscription fails")
	}
}

// TestBufferedChannel tests that subscriber channels are buffered
func TestBufferedChannel(t *testing.T) {
	pub := NewPublisher()
	topic := "test-topic"
	pub.CreateTopic(topic)

	ch, err := pub.Subscribe(topic)
	if err != nil {
		t.Fatalf("Subscribe() returned error: %v", err)
	}

	// Publish two messages quickly
	// If channel is buffered (capacity 1), first message should be buffered
	// and second should block or be buffered if there's space
	err = pub.Publish(topic, "message1")
	if err != nil {
		t.Fatalf("Publish() returned error: %v", err)
	}

	// Small delay to ensure first message is sent
	time.Sleep(10 * time.Millisecond)

	// Publish second message
	err = pub.Publish(topic, "message2")
	if err != nil {
		t.Fatalf("Publish() returned error: %v", err)
	}

	// Receive messages
	msg1 := <-ch
	msg2 := <-ch

	if msg1 != "message1" {
		t.Errorf("Expected first message 'message1', got '%s'", msg1)
	}
	if msg2 != "message2" {
		t.Errorf("Expected second message 'message2', got '%s'", msg2)
	}
}

