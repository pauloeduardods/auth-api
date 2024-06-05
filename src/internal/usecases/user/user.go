package user_usecases

import (
	"auth-api/src/internal/domain/auth"
	"auth-api/src/internal/domain/user"
	"auth-api/src/internal/events"
	"auth-api/src/pkg/logger"
)

type UseCases struct {
	Register *RegisterUserUseCase
	Update   *UpdateUserUseCase
}

func NewUseCases(userService user.UserService, authService auth.AuthService, logger logger.Logger, events events.EventDispatcher) *UseCases {
	return &UseCases{
		Register: NewRegisterUserUseCase(userService, authService, logger, events),
		Update:   NewUpdateUserUseCase(userService, logger),
	}
}
