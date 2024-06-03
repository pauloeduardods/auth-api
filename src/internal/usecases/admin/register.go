package admin_usecases

import (
	"auth-api/src/internal/domain/admin"
	"auth-api/src/internal/domain/auth"
	"auth-api/src/pkg/logger"
	"context"
)

type RegisterAdminUseCase struct {
	adminService admin.AdminService
	auth         auth.AuthService
	logger       logger.Logger
}

type RegisterAdminInput struct {
	SignupAdmin auth.CreateAdminInput
	admin.CreateAdminInput
}

func NewRegisterAdminUseCase(adminService admin.AdminService, auth auth.AuthService, logger logger.Logger) *RegisterAdminUseCase {
	return &RegisterAdminUseCase{
		adminService: adminService,
		auth:         auth,
		logger:       logger,
	}
}

func (uc *RegisterAdminUseCase) Execute(ctx context.Context, input RegisterAdminInput) (execErr error) {
	if err := input.SignupAdmin.Validate(); err != nil {
		return err
	}

	getByEmailInput := &admin.GetAdminByEmailInput{
		Email: input.CreateAdminInput.Email,
	}
	if err := getByEmailInput.Validate(); err != nil {
		return err
	}
	if exists, err := uc.adminService.GetByEmail(getByEmailInput); err != nil {
		if err != admin.ErrAdminNotFound {
			return err
		}
	} else if exists != nil {
		return admin.ErrAdminAlreadyExists
	}

	signUpOutput, err := uc.auth.CreateAdmin(ctx, input.SignupAdmin)
	if err != nil {
		return err
	}
	defer func() {
		if execErr != nil {
			if err := signUpOutput.Rollback(ctx); err != nil {
				uc.logger.Error("Error rolling back sign up: %s", err)
			}
		}
	}()

	adminId, err := admin.ParseAdminID(signUpOutput.Id)
	if err != nil {
		return err
	}

	input.CreateAdminInput.ID = adminId

	if err := input.CreateAdminInput.Validate(); err != nil {
		return err
	}

	createOut, err := uc.adminService.Create(&input.CreateAdminInput)
	if err != nil {
		return err
	}
	defer func() {
		if execErr != nil {
			if err := createOut.Rollback(ctx); err != nil {
				uc.logger.Error("Error rolling back create admin: %s", err)
			}
		}
	}()

	return nil
}
