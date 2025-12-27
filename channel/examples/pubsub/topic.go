package main

func (p *Publisher) CreateTopic(topic string) {
	p.Lock()
	defer p.Unlock()
	p.subscribers[topic] = make([]chan string, 0)
}
