package auth_usecases

import (
	"context"
	"monitoring-system/server/src/internal/domain/auth"
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
	removeMFAInput, err := auth.NewAdminRemoveMFAInput(input.Username)
	if err != nil {
		return err
	}

	return uc.auth.AdminRemoveMFA(ctx, removeMFAInput)
}
