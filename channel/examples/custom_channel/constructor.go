package main

import (
	"container/list"
	"sync"
)

type Channel[G any] struct {
	store    *list.List
	capacity int
	cond     *sync.Cond
	close    bool
}

func NewChannel[G any](capacity int) *Channel[G] {
	return &Channel[G]{
		store:    list.New(),
		capacity: capacity,
		cond:     sync.NewCond(&sync.Mutex{}),
		close:    false,
	}
}
