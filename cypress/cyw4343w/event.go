package cyw4343w

import "sync"

const (
	maxPendingEvents = 8
	maxEventWaiters  = 8
)

type Event interface {
	Close()
	TypeId() uint32
}

type AsyncEvent struct {
	Type   uint32
	Status uint32
	Reason uint32
	Auth   uint32

	Data          []byte
	EventOffset   int
	PayloadOffset int
	Handle        BufferHandle
}

func (event AsyncEvent) Close() {
	event.Handle.Close()
}

func (event AsyncEvent) TypeId() uint32 {
	return event.Type
}

func (event AsyncEvent) Payload() []byte {
	return event.Data[event.PayloadOffset:]
}

type Waiter[EventT Event] struct {
	events [maxPendingEvents]uint32
	count  int

	queue [maxPendingEvents]EventT
	head  int
	tail  int
	size  int

	closed bool
}

func (w *Waiter[EventT]) Match(eventType uint32) bool {
	for i := 0; i < w.count; i++ {
		if w.events[i] == eventType {
			return true
		}
	}
	return false
}

func (w *Waiter[EventT]) Push(event EventT) bool {
	if w.size == len(w.queue) {
		return false
	}

	w.queue[w.tail] = event
	w.tail++
	if w.tail == len(w.queue) {
		w.tail = 0
	}
	w.size++
	return true
}

func (w *Waiter[EventT]) Pop() (EventT, bool) {
	var zero EventT
	if w.size == 0 {
		return zero, false
	}

	event := w.queue[w.head]
	w.queue[w.head] = zero

	w.head++
	if w.head == len(w.queue) {
		w.head = 0
	}

	w.size--
	return event, true
}

type Dispatcher[EventT Event] struct {
	mu      sync.Mutex
	waiters [maxEventWaiters]*Waiter[EventT]
}

func (d *Dispatcher[EventT]) Watch(events ...uint32) *Waiter[EventT] {
	w := &Waiter[EventT]{}

	n := len(events)
	if n > len(w.events) {
		n = len(w.events)
	}

	for i := 0; i < n; i++ {
		w.events[i] = events[i]
	}
	w.count = n

	d.mu.Lock()
	for i := 0; i < len(d.waiters); i++ {
		if d.waiters[i] == nil {
			d.waiters[i] = w
			d.mu.Unlock()
			return w
		}
	}
	d.mu.Unlock()

	return nil
}

func (d *Dispatcher[EventT]) Unwatch(w *Waiter[EventT]) {
	if w == nil {
		return
	}

	d.mu.Lock()

	for i := 0; i < len(d.waiters); i++ {
		if d.waiters[i] == w {
			d.waiters[i] = nil
			break
		}
	}

	w.closed = true

	// Release any events still queued to this waiter.
	for {
		event, ok := w.Pop()
		if !ok {
			break
		}
		event.Close()
	}

	d.mu.Unlock()
}

func (d *Dispatcher[EventT]) Dispatch(event EventT) bool {
	d.mu.Lock()

	for i := 0; i < len(d.waiters); i++ {
		w := d.waiters[i]
		if w == nil || w.closed {
			continue
		}

		if !w.Match(event.TypeId()) {
			continue
		}

		if !w.Push(event) {
			d.mu.Unlock()
			return false
		}

		d.mu.Unlock()
		return true
	}

	d.mu.Unlock()
	return false
}
