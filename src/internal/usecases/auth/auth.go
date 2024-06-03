package auth_usecases

import (
	"auth-api/src/internal/domain/admin"
	"auth-api/src/internal/domain/auth"
	"auth-api/src/internal/domain/user"
	"auth-api/src/pkg/logger"
)

type UseCases struct {
	Login          *LoginUseCase
	AddGroup       *AddGroupUseCase
	RemoveGroup    *RemoveGroupUseCase
	RefreshToken   *RefreshTokenUseCase
	AddMFA         *AddMFAUseCase
	VerifyMFA      *VerifyMFAUseCase
	AdminRemoveMFA *AdminRemoveMFAUseCase
	RemoveMFA      *RemoveMFAUseCase
	ConfirmSignUp  *ConfirmSignUpUseCase
	GetMe          *GetMeUseCase
	ActivateMFA    *ActivateMFAUseCase
	Logout         *LogoutUseCase
	SetPassword    *SetPasswordUseCase
}

func NewUseCases(authService auth.AuthService, adminService admin.AdminService, userService user.UserService, logger logger.Logger) *UseCases {
	return &UseCases{
		Login:          NewLoginUseCase(authService),
		AddGroup:       NewAddGroupUseCase(adminService, userService, authService, logger),
		RemoveGroup:    NewRemoveGroupUseCase(authService, logger),
		RefreshToken:   NewRefreshTokenUseCase(authService),
		AddMFA:         NewAddMFAUseCase(authService),
		VerifyMFA:      NewVerifyMFAUseCase(authService),
		AdminRemoveMFA: NewAdminRemoveMFAUseCase(authService),
		RemoveMFA:      NewRemoveMFAUseCase(authService),
		ConfirmSignUp:  NewConfirmSignUpUseCase(authService),
		GetMe:          NewGetMeUseCase(authService),
		ActivateMFA:    NewActivateMFAUseCase(authService),
		Logout:         NewLogoutUseCase(authService),
		SetPassword:    NewSetPasswordUseCase(authService),
	}
}
