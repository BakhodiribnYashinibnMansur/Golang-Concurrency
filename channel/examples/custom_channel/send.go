package main

import "errors"

func (ch *Channel[G]) Send(message G) error {
	cond := ch.cond
	cond.L.Lock()
	defer cond.L.Unlock()
	if ch.close {
		return errors.New("channel is already closed")
	}
	for ch.store.Len() == ch.capacity {
		cond.Wait()
	}
	ch.store.PushBack(message)
	cond.Broadcast()
	return nil
}
