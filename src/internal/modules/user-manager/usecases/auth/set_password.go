package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"context"
)

type SetPasswordUseCase struct {
	auth auth.AuthService
}

type SetPasswordInput struct {
	auth.SetPasswordInput
}

func NewSetPasswordUseCase(auth auth.AuthService) *SetPasswordUseCase {
	return &SetPasswordUseCase{
		auth: auth,
	}
}

func (uc *SetPasswordUseCase) Execute(ctx context.Context, input SetPasswordInput) (*auth.LoginOutput, error) {
	if err := input.SetPasswordInput.Validate(); err != nil {
		return nil, err
	}

	return uc.auth.SetPassword(ctx, input.SetPasswordInput)
}
