package main

func (m *Mutex[T]) Get() T {
	responeChan := make(chan T)
	m.read <- responeChan
	return <-responeChan
}
