package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"context"
)

type RemoveMFAUseCase struct {
	auth auth.AuthService
}

type RemoveMFAInput struct {
	auth.RemoveMFAInput
}

func NewRemoveMFAUseCase(auth auth.AuthService) *RemoveMFAUseCase {
	return &RemoveMFAUseCase{
		auth: auth,
	}
}

func (uc *RemoveMFAUseCase) Execute(ctx context.Context, input RemoveMFAInput) error {
	if err := input.RemoveMFAInput.Validate(); err != nil {
		return err
	}

	return uc.auth.RemoveMFA(ctx, input.RemoveMFAInput)
}
