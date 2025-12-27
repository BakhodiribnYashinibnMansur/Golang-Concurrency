package main

import "errors"

func (ch *Channel[G]) Close() error {
	ch.cond.L.Lock()
	defer ch.cond.L.Unlock()
	if ch.close {
		return errors.New("close is already closed")
	}
	ch.close = true
	ch.cond.Broadcast()
	return nil
}
