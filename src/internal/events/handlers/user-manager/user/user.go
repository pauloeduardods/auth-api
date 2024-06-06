package user

import (
	"auth-api/src/internal/events"
	"auth-api/src/internal/modules/user-manager/domain/user"
	user_events "auth-api/src/internal/modules/user-manager/events/user"
	"auth-api/src/internal/modules/user-manager/usecases/auth"
	"auth-api/src/pkg/logger"
)

type EventsHandlers struct {
	sendConfirmationHandler events.EventHandler
}

func NewEventsHandlers(
	logger logger.Logger,
	authUsecases auth.UseCases,
) *EventsHandlers {
	return &EventsHandlers{
		sendConfirmationHandler: user_events.NewSendConfirmationHandler(logger, authUsecases),
	}
}

func (h *EventsHandlers) RegisterHandlers(dispatcher events.EventDispatcher) {
	dispatcher.Register(user.UserRegistered, h.sendConfirmationHandler)
}
