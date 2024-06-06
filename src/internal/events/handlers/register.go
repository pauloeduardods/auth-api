package events_handlers

import (
	"auth-api/src/internal/events"
	user_manager "auth-api/src/internal/events/handlers/user-manager/user"
	auth_usecases "auth-api/src/internal/modules/user-manager/usecases/auth"
	"auth-api/src/pkg/logger"
)

type EventsHandlers struct {
	userManagerHandlers *user_manager.EventsHandlers
}

func NewEventsHandlers(
	logger logger.Logger,
	authUsecases auth_usecases.UseCases,
) *EventsHandlers {
	return &EventsHandlers{
		userManagerHandlers: user_manager.NewEventsHandlers(logger, authUsecases),
	}
}

func (h *EventsHandlers) RegisterHandlers(dispatcher events.EventDispatcher) {
	h.userManagerHandlers.RegisterHandlers(dispatcher)
}
