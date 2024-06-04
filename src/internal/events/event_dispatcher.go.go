package events

import (
	"auth-api/src/internal/domain/events"
	"errors"
	"sync"
)

type EventDispatcherImpl struct {
	handlers map[events.EventType][]events.EventHandler
	mu       sync.RWMutex
}

func NewEventDispatcher() events.EventDispatcher {
	return &EventDispatcherImpl{
		handlers: make(map[events.EventType][]events.EventHandler),
	}
}

func (d *EventDispatcherImpl) Register(eventType events.EventType, handler events.EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *EventDispatcherImpl) Dispatch(event events.Event) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	handlers, ok := d.handlers[event.GetType()]
	if !ok {
		return errors.New("no handlers registered for event type")
	}

	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			return err //TODO: improve error handling with a list of errors
		}
	}
	return nil
}
