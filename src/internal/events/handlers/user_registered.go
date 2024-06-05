package events_handlers

import (
	"auth-api/src/internal/domain/code"
	"auth-api/src/internal/domain/email"
	"auth-api/src/internal/events"
	auth_usecases "auth-api/src/internal/usecases/auth"
	"auth-api/src/pkg/logger"
	"context"
	"fmt"
	"time"
)

type SendConfirmationHandler struct {
	logger       logger.Logger
	authUsecases auth_usecases.UseCases
}

func NewSendConfirmationHandler(logger logger.Logger, authUseCases auth_usecases.UseCases) *SendConfirmationHandler {
	return &SendConfirmationHandler{
		logger:       logger,
		authUsecases: authUseCases,
	}
}

func (h *SendConfirmationHandler) Handle(event events.Event) error {
	userRegisteredEvent, ok := event.(*events.UserRegisteredEvent)
	if !ok {
		return nil
	}

	if !userRegisteredEvent.NeedsVerification {
		return nil
	}

	expiresAt := time.Now().Add(10 * time.Minute)
	identifier := fmt.Sprintf("%s#%s", "CONFIRMATION_CODE", userRegisteredEvent.Email)

	sendCodeInput := auth_usecases.SendCodeInput{
		GenerateAndSaveInput: code.GenerateAndSaveInput{
			Identifier:        identifier,
			ExpiresAt:         expiresAt,
			Length:            6,
			CanContainLetters: false,
		},
		Email: email.Email{
			To:      userRegisteredEvent.Email,
			Subject: "Please confirm your email",
			Body:    "Your confirmation code is: %s",
		},
	}

	if err := h.authUsecases.SendCode.Execute(context.Background(), sendCodeInput); err != nil {
		h.logger.Error("failed to send confirmation code: %v", err)
		return err
	}

	h.logger.Info("confirmation email sent successfully")
	return nil
}
