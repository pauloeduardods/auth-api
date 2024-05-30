package auth

import (
	"context"
	"monitoring-system/server/src/domain/auth"
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
	removeMFAInput, err := auth.NewRemoveMFAInput(input.AccessToken)
	if err != nil {
		return err
	}

	return uc.auth.RemoveMFA(ctx, removeMFAInput)
}
