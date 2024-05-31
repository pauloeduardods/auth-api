package user_usecases

import (
	"context"
	"monitoring-system/server/src/internal/domain/user"
	"monitoring-system/server/src/pkg/logger"
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

	backup, err := uc.userService.Update(&input.UpdateUserInput)
	if err != nil {
		execErr = err
		return err
	}
	defer func() {
		if execErr != nil {
			if rollbackErr := uc.userService.RollbackUpdate(backup); rollbackErr != nil {
				uc.logger.Error("RollbackUpdate error: %v", rollbackErr)
			}
		}
	}()

	return nil
}
