package auth

import (
	"monitoring-system/server/src/domain/auth"
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
}

func NewUseCases(authService auth.AuthService) *UseCases {
	return &UseCases{
		Login:          NewLoginUseCase(authService),
		AddGroup:       NewAddGroupUseCase(authService),
		RemoveGroup:    NewRemoveGroupUseCase(authService),
		RefreshToken:   NewRefreshTokenUseCase(authService),
		AddMFA:         NewAddMFAUseCase(authService),
		VerifyMFA:      NewVerifyMFAUseCase(authService),
		AdminRemoveMFA: NewAdminRemoveMFAUseCase(authService),
		RemoveMFA:      NewRemoveMFAUseCase(authService),
		ConfirmSignUp:  NewConfirmSignUpUseCase(authService),
		GetMe:          NewGetMeUseCase(authService),
		ActivateMFA:    NewActivateMFAUseCase(authService),
		Logout:         NewLogoutUseCase(authService),
	}
}
