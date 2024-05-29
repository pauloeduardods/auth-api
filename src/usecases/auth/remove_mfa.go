package auth

import (
	"context"
	"monitoring-system/server/src/domain/auth"
)

type RemoveMFAUseCase struct {
	auth auth.Auth
}

type RemoveMFAInput struct {
	auth.RemoveMFAInput
}

func NewRemoveMFAUseCase(auth auth.Auth) *RemoveMFAUseCase {
	return &RemoveMFAUseCase{
		auth: auth,
	}
}

func (uc *RemoveMFAUseCase) Execute(ctx context.Context, input RemoveMFAInput) error {
	removeMFAInput, err := auth.NewRemoveMFAInput(input.Username)
	if err != nil {
		return err
	}

	return uc.auth.RemoveMFA(ctx, removeMFAInput)
}
