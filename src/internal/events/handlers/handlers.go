package events_handlers

import (
	"auth-api/src/internal/domain/code"
	"auth-api/src/internal/domain/email"
	user_events "auth-api/src/internal/events/handlers/user"
	"auth-api/src/pkg/logger"
)

type EventsHandlers struct {
	SendConfirmationHandler *user_events.SendConfirmationHandler
}

func NewEventsHandlers(
	logger logger.Logger,
	codeService code.CodeService,
	emailService email.EmailService,
) *EventsHandlers {
	return &EventsHandlers{
		SendConfirmationHandler: user_events.NewSendConfirmationHandler(logger, codeService, emailService),
	}
}
