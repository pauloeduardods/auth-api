package auth_usecases

import (
	"context"
	"monitoring-system/server/src/internal/domain/auth"
)

type ConfirmSignUpUseCase struct {
	auth auth.AuthService
}

type ConfirmSignUpInput struct {
	auth.ConfirmSignUpInput
}

func NewConfirmSignUpUseCase(auth auth.AuthService) *ConfirmSignUpUseCase {
	return &ConfirmSignUpUseCase{
		auth: auth,
	}
}

func (uc *ConfirmSignUpUseCase) Execute(ctx context.Context, input ConfirmSignUpInput) (*auth.ConfirmSignUpOutput, error) {
	if err := input.ConfirmSignUpInput.Validate(); err != nil {
		return nil, err
	}

	return uc.auth.ConfirmSignUp(ctx, input.ConfirmSignUpInput)
}
