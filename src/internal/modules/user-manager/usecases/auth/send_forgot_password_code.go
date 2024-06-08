package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/pkg/logger"
	"context"
)

type SendForgotPasswordCodeUseCase struct {
	logger logger.Logger
	auth   auth.AuthService
}

func NewSendForgotPasswordCodeUseCase(logger logger.Logger, auth auth.AuthService) *SendForgotPasswordCodeUseCase {
	return &SendForgotPasswordCodeUseCase{
		logger: logger,
		auth:   auth,
	}
}

type SendForgotPasswordCodeInput struct {
	Username string
}

func (sc *SendForgotPasswordCodeUseCase) Execute(ctx context.Context, input SendForgotPasswordCodeInput) error {
	generateAndSaveInput := auth.GenerateAndSendCodeInput{
		Username:   input.Username,
		Identifier: "FORGOT_PASSWORD_CODE",
		Subject:    "Reset your password",
		Body:       "Your reset password code is: %s",
	}

	if err := generateAndSaveInput.Validate(); err != nil {
		return err
	}

	_, err := sc.auth.GenerateAndSendCode(ctx, generateAndSaveInput)
	if err != nil {
		sc.logger.Error("failed to generate code: %v", err)
		return err
	}

	return nil
}
