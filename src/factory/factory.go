package factory

import (
	"auth-api/src/config"
	"auth-api/src/internal/domain/admin"
	"auth-api/src/internal/domain/auth"
	"auth-api/src/internal/domain/user"
	admin_repo "auth-api/src/internal/infra/admin/repository"
	admin_service "auth-api/src/internal/infra/admin/service"
	auth_service "auth-api/src/internal/infra/auth/service"
	user_repo "auth-api/src/internal/infra/user/repository"
	user_service "auth-api/src/internal/infra/user/service"
	admin_usecases "auth-api/src/internal/usecases/admin"
	auth_usecases "auth-api/src/internal/usecases/auth"
	user_usecases "auth-api/src/internal/usecases/user"
	"auth-api/src/pkg/jwt_verify"
	"auth-api/src/pkg/logger"
	"context"
	"database/sql"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type Factory struct {
	UseCases   UseCases
	Repository Repository
	Service    Service
}

type Service struct {
	Auth  auth.AuthService
	User  user.UserService
	Admin admin.AdminService
}

type Repository struct {
	User  user.UserRepository
	Admin admin.AdminRepository
}

type UseCases struct {
	Auth  *auth_usecases.UseCases
	User  *user_usecases.UseCases
	Admin *admin_usecases.UseCases
}

func newAuthService(logger logger.Logger, awsConfig *aws.Config, config config.Config) auth.AuthService {
	cognitoClient := cognitoidentityprovider.NewFromConfig(*awsConfig)
	jwtVerify := jwt_verify.NewAuth(config.Aws.Region, config.Aws.CognitoUserPoolID, logger)
	jwtVerify.CacheJWK() //TODO: Check when we need to cache the JWK and how to handle the error
	return auth_service.NewAuthService(cognitoClient, config.Aws.CognitoClientId, jwtVerify, config.Aws.CognitoUserPoolID, logger)
}

func newAuthUseCases(auth auth.AuthService) *auth_usecases.UseCases {
	return auth_usecases.NewUseCases(auth)
}

func newUserRepo(db *sql.DB, logger logger.Logger) user.UserRepository {
	return user_repo.NewUserRepository(db, logger)
}

func newUserService(userRepo user.UserRepository) user.UserService {
	return user_service.NewUserService(userRepo)
}

func newUserUseCases(userService user.UserService, authService auth.AuthService, logger logger.Logger) *user_usecases.UseCases {
	return user_usecases.NewUseCases(userService, authService, logger)
}

func newAdminRepo(db *sql.DB, logger logger.Logger) admin.AdminRepository {
	return admin_repo.NewAdminRepository(db, logger)
}

func newAdminService(adminRepo admin.AdminRepository, logger logger.Logger) admin.AdminService {
	return admin_service.NewAdminService(adminRepo, logger)
}

func newAdminUseCases(adminService admin.AdminService, authService auth.AuthService, logger logger.Logger) *admin_usecases.UseCases {
	return admin_usecases.NewUseCases(adminService, authService, logger)
}

func New(ctx context.Context, logger logger.Logger, awsConfig aws.Config, config config.Config, db *sql.DB) (*Factory, error) {
	authService := newAuthService(logger, &awsConfig, config)
	authUseCases := newAuthUseCases(authService)

	userRepo := newUserRepo(db, logger)
	userService := newUserService(userRepo)
	userUseCases := newUserUseCases(userService, authService, logger)

	adminRepo := newAdminRepo(db, logger)
	adminService := newAdminService(adminRepo, logger)
	adminUseCases := newAdminUseCases(adminService, authService, logger)

	return &Factory{
		Repository: Repository{
			User:  userRepo,
			Admin: adminRepo,
		},
		Service: Service{
			Auth:  authService,
			User:  userService,
			Admin: adminService,
		},
		UseCases: UseCases{
			Auth:  authUseCases,
			User:  userUseCases,
			Admin: adminUseCases,
		},
	}, nil
}
