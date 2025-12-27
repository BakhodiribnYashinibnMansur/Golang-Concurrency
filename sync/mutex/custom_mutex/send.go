package main

func (m *Mutex[T]) Send(value T) {
	m.write <- value
}
