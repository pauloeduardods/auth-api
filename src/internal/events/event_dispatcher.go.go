package events

import (
	"auth-api/src/pkg/logger"
	"errors"
	"sync"
)

type EventDispatcherImpl struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
	logger   logger.Logger
}

func NewEventDispatcher(logger logger.Logger) *EventDispatcherImpl {
	return &EventDispatcherImpl{
		handlers: make(map[EventType][]EventHandler),
		logger:   logger,
	}
}

func (d *EventDispatcherImpl) Register(eventType EventType, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *EventDispatcherImpl) Dispatch(event Event) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	handlers, ok := d.handlers[event.GetType()]
	if !ok {
		return errors.New("no handlers registered for event type")
	}

	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h.Handle(event); err != nil {
				d.logger.Error("error handling event: %v err: %v", event, err)
			}
		}(handler)
	}
	return nil
}
