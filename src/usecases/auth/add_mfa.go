package auth

import (
	"context"
	"monitoring-system/server/src/domain/auth"
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
	addMFAInput, err := auth.NewAddMFAInput(input.AccessToken)
	if err != nil {
		return nil, err
	}

	return uc.auth.AddMFA(ctx, addMFAInput)
}
