package cyw4343w

import "sync"

type queue struct {
	head  *queueEntry
	mutex *sync.Mutex
}

type queueEntry struct {
	next  *queueEntry
	value []byte
	id    uintptr
}

func (q *queue) Insert(id uintptr, value []byte) {
	q.mutex.Lock()
	entry := &queueEntry{
		value: value,
		next:  q.head,
		id:    id,
	}
	q.head = entry
	q.mutex.Unlock()
}

func (q *queue) Dequeue(id uintptr) ([]byte, bool) {
	q.mutex.Lock()
	var last *queueEntry
	curr := q.head
	for curr != nil {
		if curr.id == id {
			// Remove from the linked list.
			if last == nil {
				q.head = q.head.next
			} else {
				last.next = curr.next
			}

			// Return the value.
			return curr.value, true
		}

		// Advance.
		last = curr
		curr = curr.next
	}
	q.mutex.Unlock()
	return nil, false
}
