package auth_usecases

import (
	"auth-api/src/internal/domain/auth"
	"auth-api/src/internal/domain/code"
	"context"
	"fmt"
)

type ConfirmSignUpUseCase struct {
	auth auth.AuthService
	code code.CodeService
}

type ConfirmSignUpInput struct {
	auth.VerifyEmailInput
	code.VerifyCodeInput
	auth.ConfirmSignUpInput
}

func NewConfirmSignUpUseCase(auth auth.AuthService, code code.CodeService) *ConfirmSignUpUseCase {
	return &ConfirmSignUpUseCase{
		auth: auth,
		code: code,
	}
}

func (uc *ConfirmSignUpUseCase) Execute(ctx context.Context, input ConfirmSignUpInput) (*auth.ConfirmSignUpOutput, error) {
	if err := input.ConfirmSignUpInput.Validate(); err != nil {
		return nil, err
	}
	if err := uc.auth.VerifyEmail(ctx, input.VerifyEmailInput); err != nil {
		return nil, err
	}

	input.VerifyCodeInput.Identifier = fmt.Sprintf("%s#%s", "CONFIRMATION_CODE", input.ConfirmSignUpInput.Username)

	if err := input.VerifyCodeInput.Validate(); err != nil {
		return nil, err
	}

	if err := uc.code.VerifyCode(input.VerifyCodeInput); err != nil {
		return nil, err
	}

	out, err := uc.auth.ConfirmSignUp(ctx, input.ConfirmSignUpInput)
	if err != nil {
		return nil, err
	}

	if err := uc.auth.VerifyEmail(ctx, input.VerifyEmailInput); err != nil {
		return nil, err
	}

	return out, nil
}
