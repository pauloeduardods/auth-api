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

type RegisterUserInput struct {
	Email    string
	Password string
	Name     string
	Phone    *string
}

func NewRegisterUserUseCase(userService user.UserService, auth auth.AuthService, logger logger.Logger) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userService: userService,
		auth:        auth,
		logger:      logger,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (execErr error) {
	signUpInput, err := auth.NewSignUpInput(input.Email, input.Password, input.Name)
	if err != nil {
		return err
	}

	signUpOutput, err := uc.auth.SignUp(ctx, signUpInput)
	if err != nil {
		return err
	}
	defer func() {
		if execErr != nil {
			deleteUserInput, err := auth.NewDeleteUserInput(input.Email)
			if err != nil {
				uc.logger.Error("RollbackSignUp error: %v", err)
				return
			}
			if rollbackErr := uc.auth.DeleteUser(ctx, deleteUserInput); rollbackErr != nil {
				uc.logger.Error("RollbackSignUp error: %v", rollbackErr)
			}
		}
	}()

	createUserInput, err := user.NewCreateUserInput(signUpOutput.Id, input.Name, input.Email, input.Phone)
	if err != nil {
		execErr = err
		return err
	}

	err = uc.userService.Create(&createUserInput)
	if err != nil {
		execErr = err
		return err
	}
	defer func() {
		if execErr != nil {
			if rollbackErr := uc.userService.RollbackCreate(&createUserInput); rollbackErr != nil {
				uc.logger.Error("RollbackCreate error: %v", rollbackErr)
			}
		}
	}()

	return nil
}
