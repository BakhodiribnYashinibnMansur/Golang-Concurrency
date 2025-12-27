package main

type Mutex[T any] struct {
	data  T
	read  chan chan T
	write chan T
	stop  chan struct{}
}
