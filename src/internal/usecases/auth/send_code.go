package auth_usecases

import (
	"auth-api/src/internal/domain/code"
	"auth-api/src/internal/domain/email"
	"auth-api/src/pkg/logger"
	"context"
	"fmt"
)

type SendCodeUseCase struct {
	logger       logger.Logger
	codeService  code.CodeService
	emailService email.EmailService
}

func NewSendCodeUseCase(logger logger.Logger, codeService code.CodeService, emailService email.EmailService) *SendCodeUseCase {
	return &SendCodeUseCase{
		logger:       logger,
		codeService:  codeService,
		emailService: emailService,
	}
}

type SendCodeInput struct {
	code.GenerateAndSaveInput
	email.Email
}

func (sc *SendCodeUseCase) Execute(ctx context.Context, input SendCodeInput) error {
	if err := input.GenerateAndSaveInput.Validate(); err != nil {
		return err
	}

	code, err := sc.codeService.GenerateAndSave(input.GenerateAndSaveInput)
	if err != nil {
		sc.logger.Error("failed to generate code: %v", err)
		return err
	}

	email := email.Email{
		To:      input.Email.To,
		Subject: input.Email.Subject,
		Body:    fmt.Sprintf(input.Email.Body, code.Value), // Ex: Your confirmation code is: %s
	}

	if err := input.Email.Validate(); err != nil {
		return err
	}

	if err := sc.emailService.SendEmail(email); err != nil {
		sc.logger.Error("failed to send email: %v", err)
		return err
	}

	return nil
}
