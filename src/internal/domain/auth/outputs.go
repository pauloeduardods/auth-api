package auth

import (
	"context"
	"monitoring-system/server/src/pkg/app_error"
)

type LoginOutput struct {
	AccessToken  string `json:"accessToken,omitempty"`
	IdToken      string `json:"idToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	Session      string `json:"session,omitempty"`
}

type SignUpOutput struct {
	IsConfirmed bool   `json:"isConfirmed"`
	Id          string `json:"id"`
	Username    string `json:"username"`
}

func (output *SignUpOutput) Rollback(ctx context.Context, authService AuthService) error {
	deleteUserInput := DeleteUserInput{
		Username: output.Username,
	}
	if err := deleteUserInput.Validate(); err != nil {
		return app_error.NewApiError(500, "Failed to validate rollback", err.Error())
	}
	if err := authService.DeleteUser(ctx, deleteUserInput); err != nil {
		return app_error.NewApiError(500, "Failed to rollback user creation", err.Error())
	}
	return nil
}

type ConfirmSignUpOutput struct {
}

type RefreshTokenOutput struct {
	AccessToken string `json:"accessToken"`
	IdToken     string `json:"idToken"`
}

type GetMeOutput struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type CreateAdminOutput struct {
	Username string `json:"username"`
}

type AddMFAOutput struct {
	SecretCode string  `json:"secretCode"`
	Session    *string `json:"session,omitempty"`
}
