package user

import (
	"context"
	"monitoring-system/server/src/domain/user"
	"monitoring-system/server/src/pkg/logger"
)

type UpdateUserUseCase struct {
	userService user.UserService
	logger      logger.Logger
}

type UpdateUserInput struct {
	// Email *string
	Id    string
	Name  *string
	Phone *string
}

func NewUpdateUserUseCase(userService user.UserService, logger logger.Logger) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userService: userService,
		logger:      logger,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (execErr error) {
	updateUserInput, err := user.NewUpdateUserInput(input.Id, input.Name, nil, input.Phone)
	if err != nil {
		execErr = err
		return err
	}

	backup, err := uc.userService.Update(&updateUserInput)
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
