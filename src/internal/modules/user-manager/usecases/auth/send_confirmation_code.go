package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/pkg/logger"
	"context"
)

type SendConfirmationCodeUseCase struct {
	logger logger.Logger
	auth   auth.AuthService
}

func NewSendConfirmationCodeUseCase(logger logger.Logger, auth auth.AuthService) *SendConfirmationCodeUseCase {
	return &SendConfirmationCodeUseCase{
		logger: logger,
		auth:   auth,
	}
}

type SendConfirmationCodeInput struct {
	Username string
}

func (sc *SendConfirmationCodeUseCase) Execute(ctx context.Context, input SendConfirmationCodeInput) error {
	getUserInput := auth.GetUserInput{
		Username: input.Username,
	}
	if err := getUserInput.Validate(); err != nil {
		return err
	}

	getUserOutput, err := sc.auth.GetUser(ctx, getUserInput)
	if err != nil {
		return err
	}
	if getUserOutput == nil {
		return auth.ErrUserNotFound
	}

	if getUserOutput.Status == auth.Confirmed {
		return auth.ErrUserAlreadyConfirmed
	}

	generateAndSaveInput := auth.GenerateAndSendCodeInput{
		Username:   input.Username,
		Identifier: "CONFIRMATION_CODE",
		Subject:    "Please confirm your email",
		Body:       "Your confirmation code is: %s",
	}

	if err := generateAndSaveInput.Validate(); err != nil {
		return err
	}

	_, err = sc.auth.GenerateAndSendCode(ctx, generateAndSaveInput)
	if err != nil {
		sc.logger.Error("failed to generate code: %v", err)
		return err
	}

	return nil
}
