package events_handlers

import (
	"auth-api/src/internal/events"
	auth_usecases "auth-api/src/internal/usecases/auth"
	"auth-api/src/pkg/logger"
)

type EventsHandlers struct {
	sendConfirmationHandler *SendConfirmationHandler
}

func NewEventsHandlers(
	logger logger.Logger,
	authUsecases auth_usecases.UseCases,
) *EventsHandlers {
	return &EventsHandlers{
		sendConfirmationHandler: NewSendConfirmationHandler(logger, authUsecases),
	}
}

func (h *EventsHandlers) RegisterHandlers(dispatcher events.EventDispatcher) {
	dispatcher.Register(events.UserRegistered, h.sendConfirmationHandler)
}
