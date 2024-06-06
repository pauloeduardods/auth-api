package auth

import (
	"auth-api/src/internal/modules/user-manager/domain/admin"
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/internal/modules/user-manager/domain/user"
	"auth-api/src/pkg/logger"
	"context"
)

type AddGroupUseCase struct {
	adminService admin.AdminService
	userService  user.UserService
	auth         auth.AuthService
	logger       logger.Logger
}

type AddGroupInput struct {
	auth.AddGroupInput
	CreateAdminInput *admin.CreateAdminInput
	CreateUserInput  *user.CreateUserInput
}

func NewAddGroupUseCase(adminService admin.AdminService, userService user.UserService, auth auth.AuthService, logger logger.Logger) *AddGroupUseCase {
	return &AddGroupUseCase{
		adminService: adminService,
		auth:         auth,
		logger:       logger,
		userService:  userService,
	}
}

func (uc *AddGroupUseCase) Execute(ctx context.Context, input AddGroupInput) (execErr error) {
	if err := input.AddGroupInput.Validate(); err != nil {
		return err
	}

	getUserInput := auth.GetUserInput{
		Username: input.AddGroupInput.Username,
	}
	if err := getUserInput.Validate(); err != nil {
		return err
	}

	getUserOutput, err := uc.auth.GetUser(ctx, getUserInput)
	if err != nil {
		return err
	}
	if getUserOutput == nil {
		return auth.ErrUserNotFound
	}

	switch input.AddGroupInput.GroupName {
	case auth.Admin:
		getByEmailInput := &admin.GetAdminByEmailInput{
			Email: input.CreateAdminInput.Email,
		}
		if err := getByEmailInput.Validate(); err != nil {
			return err
		}

		exists, err := uc.adminService.GetByEmail(getByEmailInput)
		if err != nil && err != admin.ErrAdminNotFound {
			return err
		}

		if exists == nil {
			adminId, err := admin.ParseAdminID(getUserOutput.Id)
			if err != nil {
				return err
			}

			input.CreateAdminInput.ID = adminId

			if err := input.CreateAdminInput.Validate(); err != nil {
				return err
			}

			createOut, err := uc.adminService.Create(input.CreateAdminInput)
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
		}

	case auth.User:
		getByEmailInput := &user.GetUserByEmailInput{
			Email: input.CreateUserInput.Email,
		}

		if err := getByEmailInput.Validate(); err != nil {
			return err
		}

		exists, err := uc.userService.GetByEmail(getByEmailInput)
		if err != nil && err != user.ErrUserNotFound {
			return err
		}

		if exists == nil {
			userId, err := user.ParseUserID(getUserOutput.Id)
			if err != nil {
				return err
			}

			input.CreateUserInput.ID = userId

			if err := input.CreateUserInput.Validate(); err != nil {
				return err
			}

			createOut, err := uc.userService.Create(input.CreateUserInput)
			if err != nil {
				return err
			}
			defer func() {
				if execErr != nil {
					if err := createOut.Rollback(ctx); err != nil {
						uc.logger.Error("Error rolling back create user: %s", err)
					}
				}
			}()
		}
	default:
		return auth.ErrInvalidGroup
	}

	err = uc.auth.AddGroup(ctx, input.AddGroupInput)
	if err != nil {
		return err
	}

	err = uc.auth.AdminLogout(ctx, auth.AdminLogoutInput{ // Force logout user to update token and update token with new group
		Username: input.AddGroupInput.Username,
	})
	if err != nil {
		uc.logger.Error("Error admin logging out: %s", err) //TODO: check if need to return error
	}

	return nil
}
