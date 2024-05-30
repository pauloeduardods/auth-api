package auth

import (
	"monitoring-system/server/src/domain/auth"
)

type UseCases struct {
	Login         *LoginUseCase
	AddGroup      *AddGroupUseCase
	RemoveGroup   *RemoveGroupUseCase
	RefreshToken  *RefreshTokenUseCase
	AddMFA        *AddMFAUseCase
	VerifyMFA     *VerifyMFAUseCase
	RemoveMFA     *RemoveMFAUseCase
	ConfirmSignUp *ConfirmSignUpUseCase
	GetMe         *GetMeUseCase
}

func NewUseCases(authService auth.AuthService) *UseCases {
	return &UseCases{
		Login:         NewLoginUseCase(authService),
		AddGroup:      NewAddGroupUseCase(authService),
		RemoveGroup:   NewRemoveGroupUseCase(authService),
		RefreshToken:  NewRefreshTokenUseCase(authService),
		AddMFA:        NewAddMFAUseCase(authService),
		VerifyMFA:     NewVerifyMFAUseCase(authService),
		RemoveMFA:     NewRemoveMFAUseCase(authService),
		ConfirmSignUp: NewConfirmSignUpUseCase(authService),
		GetMe:         NewGetMeUseCase(authService),
	}
}
