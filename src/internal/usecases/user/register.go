package user_usecases

import (
	"auth-api/src/internal/domain/auth"
	"auth-api/src/internal/domain/user"
	"auth-api/src/pkg/logger"
	"context"
)

type RegisterUserUseCase struct {
	userService user.UserService
	auth        auth.AuthService
	logger      logger.Logger
}

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

	getByEmailInput := &user.GetUserByEmailInput{
		Email: input.CreateUserInput.Email,
	}
	if err := getByEmailInput.Validate(); err != nil {
		return err
	}
	if exists, err := uc.userService.GetByEmail(getByEmailInput); err != nil {
		if err != user.ErrUserNotFound {
			return err
		}
	} else if exists != nil {
		return user.ErrUserAlreadyExists
	}

	signUpOutput, err := uc.auth.SignUp(ctx, input.SignUpInput)
	if err != nil {
		return err
	}
	defer func() {
		if execErr != nil {
			if err := signUpOutput.Rollback(ctx, uc.auth); err != nil {
				uc.logger.Error("Error rolling back sign up: %s", err)
			}
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
