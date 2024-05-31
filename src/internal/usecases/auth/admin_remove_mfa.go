package auth_usecases

import (
	"auth-api/src/internal/domain/auth"
	"context"
)

type AdminRemoveMFAUseCase struct {
	auth auth.AuthService
}

type AdminRemoveMFAInput struct {
	auth.AdminRemoveMFAInput
}

func NewAdminRemoveMFAUseCase(auth auth.AuthService) *AdminRemoveMFAUseCase {
	return &AdminRemoveMFAUseCase{
		auth: auth,
	}
}

func (uc *AdminRemoveMFAUseCase) Execute(ctx context.Context, input AdminRemoveMFAInput) error {
	if err := input.AdminRemoveMFAInput.Validate(); err != nil {
		return err
	}

	return uc.auth.AdminRemoveMFA(ctx, input.AdminRemoveMFAInput)
}
