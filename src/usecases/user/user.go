package user

import (
	"monitoring-system/server/src/domain/auth"
	"monitoring-system/server/src/domain/user"
	"monitoring-system/server/src/pkg/logger"
)

type UseCases struct {
	Register *RegisterUserUseCase
	Update   *UpdateUserUseCase
}

func NewUseCases(userService user.UserService, authService auth.AuthService, logger logger.Logger) *UseCases {
	return &UseCases{
		Register: NewRegisterUserUseCase(userService, authService, logger),
		Update:   NewUpdateUserUseCase(userService, logger),
	}
}
