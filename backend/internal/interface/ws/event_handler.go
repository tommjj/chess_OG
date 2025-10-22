// EventHandler manages WebSocket event handlers.
// It allows registering, retrieving, unregistering, and clearing event handlers.
package ws

import (
	"sync"
)

type HandleFunc func(ctx *Context)

type EventHandler struct {
	handlers map[string]HandleFunc
	mx       sync.RWMutex
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		handlers: make(map[string]HandleFunc),
	}
}

func (eh *EventHandler) Register(event string, handler HandleFunc) {
	eh.mx.Lock()
	defer eh.mx.Unlock()

	eh.handlers[event] = handler
}

func (eh *EventHandler) Get(event string) (HandleFunc, bool) {
	eh.mx.RLock()
	defer eh.mx.RUnlock()

	h, ok := eh.handlers[event]
	return h, ok
}

func (eh *EventHandler) MustGet(event string) HandleFunc {
	h, ok := eh.Get(event)
	if !ok {
		panic("MustGet: event handler not found for event " + event)
	}
	return h
}

func (eh *EventHandler) Unregister(event string) {
	eh.mx.Lock()
	defer eh.mx.Unlock()

	delete(eh.handlers, event)
}

func (eh *EventHandler) Clear() {
	eh.mx.Lock()
	defer eh.mx.Unlock()

	eh.handlers = make(map[string]HandleFunc)
}
