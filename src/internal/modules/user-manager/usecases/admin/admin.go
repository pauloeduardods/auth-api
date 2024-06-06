package admin

import (
	"auth-api/src/internal/modules/user-manager/domain/admin"
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/pkg/logger"
)

type UseCases struct {
	Register *RegisterAdminUseCase
	Update   *UpdateAdminUseCase
}

func NewUseCases(adminService admin.AdminService, authService auth.AuthService, logger logger.Logger) *UseCases {
	return &UseCases{
		Register: NewRegisterAdminUseCase(adminService, authService, logger),
		Update:   NewUpdateAdminUseCase(adminService, logger),
	}
}
