package auth_usecases

import (
	"context"
	"monitoring-system/server/src/internal/domain/auth"
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
	logoutInput, err := auth.NewLogoutInput(input.AccessToken)
	if err != nil {
		return err
	}

	return uc.auth.Logout(ctx, logoutInput)
}
