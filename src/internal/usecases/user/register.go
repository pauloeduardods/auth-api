package user_usecases

import (
	"context"
	"monitoring-system/server/src/internal/domain/auth"
	"monitoring-system/server/src/internal/domain/user"
	"monitoring-system/server/src/pkg/logger"
)

type RegisterUserUseCase struct {
	userService user.UserService
	auth        auth.AuthService
	logger      logger.Logger
}

// type RegisterUserInput struct {
// 	Email    string
// 	Password string
// 	Name     string
// 	Phone    *string
// }

type RegisterUserInput struct {
	auth.SignUpInput
	user.CreateUserInput
}

func NewRegisterUserUseCase(userService user.UserService, auth auth.AuthService, logger logger.Logger) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userService: userService,
		auth:        auth,
		logger:      logger,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (execErr error) {
	if err := input.SignUpInput.Validate(); err != nil {
		return err
	}

	signUpOutput, err := uc.auth.SignUp(ctx, input.SignUpInput)
	if err != nil {
		return err
	}
	defer func() {
		if execErr != nil {
			signUpOutput.Rollback(ctx, uc.auth)
		}
	}()

	userId, err := user.ParseUserID(signUpOutput.Id)
	if err != nil {
		return err
	}

	input.CreateUserInput.ID = userId

	if err := input.CreateUserInput.Validate(); err != nil {
		return err
	}

	createOut, err := uc.userService.Create(&input.CreateUserInput)
	if err != nil {
		return err
	}
	defer func() {
		if execErr != nil {
			if err := createOut.Rollback(ctx, uc.userService); err != nil {
				uc.logger.Error("Error rolling back create user: %s", err)
			}
		}
	}()

	return nil
}
