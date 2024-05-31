package user_usecases

import (
	"auth-api/src/internal/domain/user"
	"auth-api/src/pkg/logger"
	"context"
)

type UpdateUserUseCase struct {
	userService user.UserService
	logger      logger.Logger
}

type UpdateUserInput struct {
	user.UpdateUserInput
}

func NewUpdateUserUseCase(userService user.UserService, logger logger.Logger) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userService: userService,
		logger:      logger,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (execErr error) {
	if err := input.UpdateUserInput.Validate(); err != nil {
		return err
	}

	updateOut, err := uc.userService.Update(&input.UpdateUserInput)
	if err != nil {
		return err
	}
	defer func() {
		if execErr != nil {
			if err := updateOut.Rollback(ctx, uc.userService); err != nil {
				uc.logger.Error("Error rolling back update user: %s", err)
			}
		}
	}()

	return nil
}
