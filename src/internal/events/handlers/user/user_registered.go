package user_events

import (
	"auth-api/src/internal/domain/code"
	"auth-api/src/internal/domain/email"
	"auth-api/src/internal/domain/events"
	"auth-api/src/internal/domain/user"
	"auth-api/src/pkg/logger"
	"fmt"
	"time"
)

type SendConfirmationHandler struct {
	logger       logger.Logger
	codeService  code.CodeService
	emailService email.EmailService
}

func NewSendConfirmationHandler(logger logger.Logger, codeService code.CodeService, emailService email.EmailService) *SendConfirmationHandler {
	return &SendConfirmationHandler{
		logger:       logger,
		codeService:  codeService,
		emailService: emailService,
	}
}

func (h *SendConfirmationHandler) Handle(event events.Event) error {
	userRegisteredEvent, ok := event.(*user.UserRegisteredEvent)
	if !ok {
		return nil
	}

	if !userRegisteredEvent.NeedsVerification {
		return nil
	}

	expiresAt := time.Now().Add(10 * time.Minute)
	identifier := fmt.Sprintf("%s#%s", "CONFIRMATION_CODE", userRegisteredEvent.Email)

	generatedCodeInput := code.GenerateAndSaveInput{
		Identifier:        identifier,
		ExpiresAt:         expiresAt,
		Length:            6,
		CanContainLetters: false,
	}

	code, err := h.codeService.GenerateAndSave(generatedCodeInput)
	if err != nil {
		h.logger.Error("failed to generate code: %v", err)
		return err
	}

	email := email.Email{
		To:      userRegisteredEvent.Email,
		Subject: "Please confirm your email",
		Body:    "Your confirmation code is: " + code.Value,
	}

	if err := h.emailService.SendEmail(email); err != nil {
		h.logger.Error("failed to send email: %v", err)
		return err
	}

	h.logger.Info("confirmation email sent successfully")
	return nil
}
