package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"context"
)

type ConfirmSignUpUseCase struct {
	auth auth.AuthService
}

type ConfirmSignUpInput struct {
	Username string
	Code     string
}

func NewConfirmSignUpUseCase(auth auth.AuthService) *ConfirmSignUpUseCase {
	return &ConfirmSignUpUseCase{
		auth: auth,
	}
}

func (uc *ConfirmSignUpUseCase) Execute(ctx context.Context, input ConfirmSignUpInput) (*auth.ConfirmSignUpOutput, error) {
	verifyCodeInput := auth.VerifyCodeInput{
		Username:   input.Username,
		Code:       input.Code,
		Identifier: "CONFIRMATION_CODE",
	}
	if err := verifyCodeInput.Validate(); err != nil {
		return nil, err
	}

	confirmSignUpInput := auth.ConfirmSignUpInput{
		Username: input.Username,
	}
	if err := confirmSignUpInput.Validate(); err != nil {
		return nil, err
	}

	verifyEmailInput := auth.VerifyEmailInput{
		Username: input.Username,
	}
	if err := verifyEmailInput.Validate(); err != nil {
		return nil, err
	}

	if err := uc.auth.VerifyCode(ctx, verifyCodeInput); err != nil {
		return nil, err
	}

	out, err := uc.auth.ConfirmSignUp(ctx, confirmSignUpInput)
	if err != nil {
		return nil, err
	}

	if err := uc.auth.VerifyEmail(ctx, verifyEmailInput); err != nil {
		return nil, err
	}

	return out, nil
}
