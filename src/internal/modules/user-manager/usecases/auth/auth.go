package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/admin"
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/internal/modules/user-manager/domain/user"
	"auth-api/src/pkg/logger"
)

type UseCases struct {
	Login                  *LoginUseCase
	AddGroup               *AddGroupUseCase
	RemoveGroup            *RemoveGroupUseCase
	RefreshToken           *RefreshTokenUseCase
	AddMFA                 *AddMFAUseCase
	VerifyMFA              *VerifyMFAUseCase
	AdminRemoveMFA         *AdminRemoveMFAUseCase
	RemoveMFA              *RemoveMFAUseCase
	ConfirmSignUp          *ConfirmSignUpUseCase
	GetMe                  *GetMeUseCase
	ActivateMFA            *ActivateMFAUseCase
	Logout                 *LogoutUseCase
	SetPassword            *SetPasswordUseCase
	SendConfirmationCode   *SendConfirmationCodeUseCase
	ChangePassword         *ChangePasswordUseCase
	ResetPassword          *ResetPasswordUseCase
	SendForgotPasswordCode *SendForgotPasswordCodeUseCase
}

func NewUseCases(authService auth.AuthService, adminService admin.AdminService, userService user.UserService, logger logger.Logger) *UseCases {
	return &UseCases{
		Login:                  NewLoginUseCase(authService),
		AddGroup:               NewAddGroupUseCase(adminService, userService, authService, logger),
		RemoveGroup:            NewRemoveGroupUseCase(authService, logger),
		RefreshToken:           NewRefreshTokenUseCase(authService),
		AddMFA:                 NewAddMFAUseCase(authService),
		VerifyMFA:              NewVerifyMFAUseCase(authService),
		AdminRemoveMFA:         NewAdminRemoveMFAUseCase(authService),
		RemoveMFA:              NewRemoveMFAUseCase(authService),
		ConfirmSignUp:          NewConfirmSignUpUseCase(authService),
		GetMe:                  NewGetMeUseCase(authService),
		ActivateMFA:            NewActivateMFAUseCase(authService),
		Logout:                 NewLogoutUseCase(authService),
		SetPassword:            NewSetPasswordUseCase(authService),
		SendConfirmationCode:   NewSendConfirmationCodeUseCase(logger, authService),
		ChangePassword:         NewChangePasswordUseCase(authService),
		ResetPassword:          NewResetPasswordUseCase(authService),
		SendForgotPasswordCode: NewSendForgotPasswordCodeUseCase(logger, authService),
	}
}
