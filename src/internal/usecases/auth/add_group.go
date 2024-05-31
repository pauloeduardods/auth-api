package auth_usecases

import (
	"auth-api/src/internal/domain/auth"
	"context"
)

type AddGroupUseCase struct {
	auth auth.AuthService
}

type AddGroupInput struct {
	auth.AddGroupInput
}

func NewAddGroupUseCase(auth auth.AuthService) *AddGroupUseCase {
	return &AddGroupUseCase{
		auth: auth,
	}
}

func (uc *AddGroupUseCase) Execute(ctx context.Context, input AddGroupInput) error {
	if err := input.AddGroupInput.Validate(); err != nil {
		return err
	}

	return uc.auth.AddGroup(ctx, input.AddGroupInput)
}
