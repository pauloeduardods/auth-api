package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"context"
)

type AddMFAUseCase struct {
	auth auth.AuthService
}

type AddMFAInput struct {
	auth.AddMFAInput
}

func NewAddMFAUseCase(auth auth.AuthService) *AddMFAUseCase {
	return &AddMFAUseCase{
		auth: auth,
	}
}

func (uc *AddMFAUseCase) Execute(ctx context.Context, input AddMFAInput) (*auth.AddMFAOutput, error) {
	if err := input.AddMFAInput.Validate(); err != nil {
		return nil, err
	}

	return uc.auth.AddMFA(ctx, input.AddMFAInput)
}
