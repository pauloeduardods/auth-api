package auth_usecases

import (
	"auth-api/src/internal/domain/auth"
	"context"
)

type LogoutUseCase struct {
	auth auth.AuthService
}

type LogoutInput struct {
	auth.LogoutInput
}

func NewLogoutUseCase(auth auth.AuthService) *LogoutUseCase {
	return &LogoutUseCase{
		auth: auth,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) error {
	if err := input.LogoutInput.Validate(); err != nil {
		return err
	}

	return uc.auth.Logout(ctx, input.LogoutInput)
}
