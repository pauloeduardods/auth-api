package user

import (
	"auth-api/src/internal/events"
	"auth-api/src/internal/modules/user-manager/domain/user"
	"auth-api/src/internal/modules/user-manager/usecases/auth"
	"auth-api/src/pkg/logger"
	"context"
)

type SendConfirmationHandler struct {
	logger       logger.Logger
	authUsecases auth.UseCases
}

func NewSendConfirmationHandler(logger logger.Logger, authUseCases auth.UseCases) events.EventHandler {
	return &SendConfirmationHandler{
		logger:       logger,
		authUsecases: authUseCases,
	}
}

func (h *SendConfirmationHandler) Handle(event events.Event) error {
	userRegisteredEvent, ok := event.(*user.UserRegisteredEvent)
	if !ok {
		return nil
	}

	if err := userRegisteredEvent.Validate(); err != nil {
		return err
	}

	if !userRegisteredEvent.NeedsVerification {
		return nil
	}

	if err := h.authUsecases.SendConfirmationCodeUseCase.Execute(context.TODO(), auth.SendConfirmationCodeInput{
		Username: userRegisteredEvent.Email,
	}); err != nil {
		h.logger.Error("failed to send confirmation code: %v", err)
		return err
	}

	return nil
}
