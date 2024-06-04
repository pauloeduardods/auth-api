package register

import (
	"auth-api/src/internal/domain/events"
	events_handlers "auth-api/src/internal/events/handlers"
)

type Register struct {
	eventDispatcher events.EventDispatcher
	handlers        events_handlers.EventsHandlers
}

func NewRegister(eventDispatcher events.EventDispatcher, handlers events_handlers.EventsHandlers) *Register {
	return &Register{
		eventDispatcher: eventDispatcher,
		handlers:        handlers,
	}
}

func (r *Register) RegisterHandlers() {
	r.eventDispatcher.Register(events.UserRegistered, r.handlers.SendConfirmationHandler)
}
