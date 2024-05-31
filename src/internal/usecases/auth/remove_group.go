package auth_usecases

import (
	"auth-api/src/internal/domain/auth"
	"context"
)

type RemoveGroupUseCase struct {
	auth auth.AuthService
}

type RemoveGroupInput struct {
	auth.RemoveGroupInput
}

func NewRemoveGroupUseCase(auth auth.AuthService) *RemoveGroupUseCase {
	return &RemoveGroupUseCase{
		auth: auth,
	}
}

func (uc *RemoveGroupUseCase) Execute(ctx context.Context, input RemoveGroupInput) error {
	if err := input.RemoveGroupInput.Validate(); err != nil {
		return err
	}

	return uc.auth.RemoveGroup(ctx, input.RemoveGroupInput)
}
