package auth_usecases

import (
	"context"
	"monitoring-system/server/src/internal/domain/auth"
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
