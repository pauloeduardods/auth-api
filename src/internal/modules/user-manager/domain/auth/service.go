package auth

import (
	"auth-api/src/pkg/app_error"
	"context"
)

var (
	ErrUserNotFound = app_error.NewApiError(404, "User not found")
)

type AuthService interface {
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
	SignUp(ctx context.Context, input SignUpInput) (*SignUpOutput, error)
	DeleteUser(ctx context.Context, input DeleteUserInput) error
	ConfirmSignUp(ctx context.Context, input ConfirmSignUpInput) (*ConfirmSignUpOutput, error)
	GetMe(ctx context.Context, input GetMeInput) (*GetMeOutput, error)
	ValidateToken(ctx context.Context, token string) (*Claims, error)
	AddGroup(ctx context.Context, input AddGroupInput) error
	RemoveGroup(ctx context.Context, input RemoveGroupInput) error
	RefreshToken(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error)
	CreateAdmin(ctx context.Context, input CreateAdminInput) (*CreateAdminOutput, error)
	AddMFA(ctx context.Context, input AddMFAInput) (*AddMFAOutput, error)
	ActivateMFA(ctx context.Context, input ActivateMFAInput) error
	VerifyMFA(ctx context.Context, input VerifyMFAInput) (*LoginOutput, error)
	AdminRemoveMFA(ctx context.Context, input AdminRemoveMFAInput) error
	RemoveMFA(ctx context.Context, input RemoveMFAInput) error
	Logout(ctx context.Context, input LogoutInput) error
	SetPassword(ctx context.Context, input SetPasswordInput) (*LoginOutput, error)
	GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error)
	AdminLogout(ctx context.Context, input AdminLogoutInput) error
	VerifyEmail(ctx context.Context, input VerifyEmailInput) error
	GenerateAndSendCode(ctx context.Context, input GenerateAndSendCodeInput) (*GenerateAndSendCodeOutput, error)
	VerifyCode(ctx context.Context, input VerifyCodeInput) error
}
