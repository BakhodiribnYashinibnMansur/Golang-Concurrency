package main

func NewMutex[T any]() *Mutex[T] {
	m := &Mutex[T]{
		read:  make(chan chan T),
		write: make(chan T),
		stop:  make(chan struct{}),
	}
	go func() {
		for {
			select {
			case responeChan := <-m.read:
				responeChan <- m.data
			case value := <-m.write:
				m.data = value
			case <-m.stop:
				return
			}
		}
	}()
	return m
}

func NewMutexWithValue[T any](value T) *Mutex[T] {
	m := NewMutex[T]()
	m.data = value
	return m
}
