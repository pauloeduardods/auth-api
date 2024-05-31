package auth_usecases

import (
	"context"
	"monitoring-system/server/src/internal/domain/auth"
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
