package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"context"
)

type ChangePasswordUseCase struct {
	auth auth.AuthService
}

type ChangePasswordInput struct {
	AccessToken string
	OldPassword string
	NewPassword string
}

func NewChangePasswordUseCase(auth auth.AuthService) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		auth: auth,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) error {
	changePasswordInput := auth.ChangePasswordInput{
		AccessToken: input.AccessToken,
		OldPassword: input.OldPassword,
		NewPassword: input.NewPassword,
	}
	if err := changePasswordInput.Validate(); err != nil {
		return err
	}

	if err := uc.auth.ChangePassword(ctx, changePasswordInput); err != nil {
		return err
	}

	return nil
}
