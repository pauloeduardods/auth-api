package factory

import (
	"context"
	"database/sql"
	"monitoring-system/server/src/config"
	"monitoring-system/server/src/internal/domain/auth"
	"monitoring-system/server/src/internal/domain/user"
	auth_service "monitoring-system/server/src/internal/infra/auth/service"
	user_repo "monitoring-system/server/src/internal/infra/user/repository"
	user_service "monitoring-system/server/src/internal/infra/user/service"
	auth_usecases "monitoring-system/server/src/internal/usecases/auth"
	user_usecases "monitoring-system/server/src/internal/usecases/user"
	"monitoring-system/server/src/pkg/jwt_verify"
	"monitoring-system/server/src/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type Factory struct {
	UseCases   UseCases
	Repository Repository
	Service    Service
}

type Service struct {
	Auth auth.AuthService
	User user.UserService
}

type Repository struct {
	User user.UserRepository
}

type UseCases struct {
	Auth *auth_usecases.UseCases
	User *user_usecases.UseCases
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

func New(ctx context.Context, logger logger.Logger, awsConfig aws.Config, config config.Config, db *sql.DB) (*Factory, error) {
	authService := newAuthService(logger, &awsConfig, config)
	authUseCases := newAuthUseCases(authService)

	userRepo := newUserRepo(db, logger)
	userService := newUserService(userRepo)
	userUseCases := newUserUseCases(userService, authService, logger)

	return &Factory{
		Repository: Repository{
			User: userRepo,
		},
		Service: Service{
			Auth: authService,
			User: userService,
		},
		UseCases: UseCases{
			Auth: authUseCases,
			User: userUseCases,
		},
	}, nil
}
