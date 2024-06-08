package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"context"
)

type ResetPasswordUseCase struct {
	auth auth.AuthService
}

type ResetPasswordInput struct {
	Username    string
	Code        string
	NewPassword string
}

func NewResetPasswordUseCase(auth auth.AuthService) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		auth: auth,
	}
}

func (uc *ResetPasswordUseCase) Execute(ctx context.Context, input ResetPasswordInput) error {
	verifyCodeInput := auth.VerifyCodeInput{
		Username:   input.Username,
		Code:       input.Code,
		Identifier: "FORGOT_PASSWORD_CODE",
	}
	if err := verifyCodeInput.Validate(); err != nil {
		return err
	}

	changeForgotInput := auth.ChangeForgotPasswordInput{
		Username:    input.Username,
		NewPassword: input.NewPassword,
	}
	if err := changeForgotInput.Validate(); err != nil {
		return err
	}

	if err := uc.auth.VerifyCode(ctx, verifyCodeInput); err != nil {
		return err
	}

	if err := uc.auth.ChangeForgotPassword(ctx, changeForgotInput); err != nil {
		return err
	}

	return nil
}
