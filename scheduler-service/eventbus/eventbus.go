package eventbus

import (
	"log"
	"sync"
)

type EventBus[T any] struct {
	mutexLock   sync.RWMutex
	subscribers map[int]chan T
	nextId      int
	bufferSize  int
}

func NewEventBus[T any](bufferSize int) *EventBus[T] {
	return &EventBus[T]{
		subscribers: make(map[int]chan T),
		bufferSize:  bufferSize,
	}
}

func (b *EventBus[T]) Subscribe() (id int, ch <-chan T) {
	b.mutexLock.Lock()
	defer b.mutexLock.Unlock()

	id = b.nextId
	b.nextId++

	c := make(chan T, b.bufferSize)
	b.subscribers[id] = c
	return id, c
}

func (b *EventBus[T]) Unsubscribe(id int) {
	b.mutexLock.Lock()
	defer b.mutexLock.Unlock()

	if ch, ok := b.subscribers[id]; ok {
		close(ch)
		delete(b.subscribers, id)
	}
}

func (b *EventBus[T]) Publish(event T) {
	b.mutexLock.RLock()
	defer b.mutexLock.RUnlock()

	for _, ch := range b.subscribers {
		select {
		case ch <- event:
			// delivered
		default:
			log.Printf("event bus buffer is full: dropping event, subscriber too slow")
			// Can implement a retry mechanism later
		}
	}
}
