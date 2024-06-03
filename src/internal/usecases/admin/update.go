package admin_usecases

import (
	"auth-api/src/internal/domain/admin"
	"auth-api/src/pkg/logger"
	"context"
)

type UpdateAdminUseCase struct {
	adminService admin.AdminService
	logger       logger.Logger
}

type UpdateAdminInput struct {
	admin.UpdateAdminInput
}

func NewUpdateAdminUseCase(adminService admin.AdminService, logger logger.Logger) *UpdateAdminUseCase {
	return &UpdateAdminUseCase{
		adminService: adminService,
		logger:       logger,
	}
}

func (uc *UpdateAdminUseCase) Execute(ctx context.Context, input UpdateAdminInput) (execErr error) {
	if err := input.UpdateAdminInput.Validate(); err != nil {
		return err
	}

	updateOut, err := uc.adminService.Update(&input.UpdateAdminInput)
	if err != nil {
		return err
	}
	defer func() {
		if execErr != nil {
			if err := updateOut.Rollback(ctx); err != nil {
				uc.logger.Error("Error rolling back update admin: %s", err)
			}
		}
	}()

	return nil
}
