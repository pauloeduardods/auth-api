package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/pkg/logger"
	"context"
)

type RemoveGroupUseCase struct {
	auth   auth.AuthService
	logger logger.Logger
}

type RemoveGroupInput struct {
	auth.RemoveGroupInput
}

func NewRemoveGroupUseCase(auth auth.AuthService, logger logger.Logger) *RemoveGroupUseCase {
	return &RemoveGroupUseCase{
		auth:   auth,
		logger: logger,
	}
}

func (uc *RemoveGroupUseCase) Execute(ctx context.Context, input RemoveGroupInput) error {
	if err := input.RemoveGroupInput.Validate(); err != nil {
		return err
	}

	if err := uc.auth.RemoveGroup(ctx, input.RemoveGroupInput); err != nil {
		return err
	}

	if err := uc.auth.AdminLogout(ctx, auth.AdminLogoutInput{
		Username: input.RemoveGroupInput.Username,
	}); err != nil {
		uc.logger.Error("Failed to admin logout", err)
	}

	return nil
}
