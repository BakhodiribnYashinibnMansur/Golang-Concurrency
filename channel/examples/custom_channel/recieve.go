package main

func (ch *Channel[G]) Receive() (message G, ok bool) {
	cond := ch.cond

	cond.L.Lock()
	defer cond.L.Unlock()

	if ch.close {
		return message, ok
	}
	ch.capacity++
	cond.Broadcast()

	for ch.store.Len() == 0 {
		cond.Wait()
	}

	ch.capacity--
	item := ch.store.Front()
	ch.store.Remove(item)
	message = item.Value.(G)
	cond.Broadcast()
	return message, true
}
