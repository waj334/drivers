package cyw4343w

import "sync"

type queue[ElementT any] struct {
	head  *queueEntry[ElementT]
	mutex *sync.Mutex
}

type queueEntry[ElementT any] struct {
	next  *queueEntry[ElementT]
	value ElementT
	id    uint32
}

func (q *queue[ElementT]) Insert(id uint32, value ElementT) {
	q.mutex.Lock()
	entry := &queueEntry[ElementT]{
		value: value,
		next:  q.head,
		id:    id,
	}
	q.head = entry
	q.mutex.Unlock()
}

func (q *queue[ElementT]) Dequeue(id uint32) (ElementT, bool) {
	q.mutex.Lock()
	var last *queueEntry[ElementT]
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
			q.mutex.Unlock()
			return curr.value, true
		}

		// Advance.
		last = curr
		curr = curr.next
	}
	q.mutex.Unlock()

	var zero ElementT
	return zero, false
}
