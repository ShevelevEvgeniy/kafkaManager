package event_dispatcher

import (
	"context"
	"sync"
)

type Message struct {
	Key   []byte
	Value []byte
}

type Listener interface {
	Handle(ctx context.Context, Message Message)
}

type EventDispatcherInterface interface {
	Subscribe(eventType string, listener Listener)
	Unsubscribe(eventType string, listener Listener)
	Trigger(event string, Message Message, ctx context.Context)
}

type EventDispatcher struct {
	listeners map[string][]Listener
}

func New() *EventDispatcher {
	return &EventDispatcher{
		listeners: make(map[string][]Listener),
	}
}

func (e *EventDispatcher) Subscribe(eventType string, listener Listener) {
	mu := sync.RWMutex{}
	mu.Lock()
	defer mu.Unlock()

	e.listeners[eventType] = append(e.listeners[eventType], listener)
}

func (e *EventDispatcher) Unsubscribe(eventType string, listener Listener) {
	mu := sync.RWMutex{}
	mu.Lock()
	defer mu.Unlock()

	listeners := e.listeners[eventType]
	for i, l := range listeners {
		if l == listener {
			e.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
}

func (e *EventDispatcher) Trigger(event string, msg Message, ctx context.Context) {
	mu := sync.RWMutex{}
	mu.Lock()
	defer mu.Unlock()

	for _, listener := range e.listeners[event] {
		listener.Handle(ctx, msg)
	}
}
