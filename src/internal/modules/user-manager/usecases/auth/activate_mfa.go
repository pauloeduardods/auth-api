package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"context"
)

type ActivateMFAUseCase struct {
	auth auth.AuthService
}

type ActivateMFAInput struct {
	auth.ActivateMFAInput
}

func NewActivateMFAUseCase(auth auth.AuthService) *ActivateMFAUseCase {
	return &ActivateMFAUseCase{
		auth: auth,
	}
}

func (uc *ActivateMFAUseCase) Execute(ctx context.Context, input ActivateMFAInput) error {
	if err := input.ActivateMFAInput.Validate(); err != nil {
		return err
	}

	return uc.auth.ActivateMFA(ctx, input.ActivateMFAInput)
}
