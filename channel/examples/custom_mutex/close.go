package main

func (m *Mutex[T]) Close() {
	m.stop <- struct{}{}
}
