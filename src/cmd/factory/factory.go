package factory

import (
	"context"
	"database/sql"
	"monitoring-system/server/src/config"
	authDomain "monitoring-system/server/src/domain/auth"
	userDomain "monitoring-system/server/src/domain/user"

	auth_client "monitoring-system/server/src/internal/auth/client"
	user_repo "monitoring-system/server/src/internal/user/repository"
	"monitoring-system/server/src/pkg/jwt_verify"
	"monitoring-system/server/src/pkg/logger"
	authUseCases "monitoring-system/server/src/usecases/auth"
	userUseCases "monitoring-system/server/src/usecases/user"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type Factory struct {
	Domain   Domain
	internal Internal
	UseCases UseCases
}

type Domain struct {
	Auth authDomain.AuthService
	User userDomain.UserService
}

type Internal struct {
	AuthClient authDomain.AuthClient
	UserRepo   userDomain.UserRepository
}

type UseCases struct {
	Auth *authUseCases.UseCases
	User *userUseCases.UseCases
}

func newDomainAuth(authClient authDomain.AuthClient) authDomain.AuthService {
	return authDomain.NewAuthService(authClient)
}

func newAuthClient(logger logger.Logger, awsConfig *aws.Config, config config.Config) authDomain.AuthClient {
	cognitoClient := cognitoidentityprovider.NewFromConfig(*awsConfig)
	jwtVerify := jwt_verify.NewAuth(config.Aws.Region, config.Aws.CognitoUserPoolID, logger)
	jwtVerify.CacheJWK() //TODO: Check when we need to cache the JWK and how to handle the error
	return auth_client.NewAuthClient(cognitoClient, config.Aws.CognitoClientId, jwtVerify, config.Aws.CognitoUserPoolID, logger)
}

func newAuthUseCases(auth authDomain.AuthService) *authUseCases.UseCases {
	return authUseCases.NewUseCases(auth)
}

func newUserRepo(db *sql.DB, logger logger.Logger) userDomain.UserRepository {
	return user_repo.NewUserRepository(db, logger)
}

func newUserService(userRepo userDomain.UserRepository) userDomain.UserService {
	return userDomain.NewUserService(userRepo)
}

func newUserUseCases(userService userDomain.UserService, authService authDomain.AuthService, logger logger.Logger) *userUseCases.UseCases {
	return userUseCases.NewUseCases(userService, authService, logger)
}

func New(ctx context.Context, logger logger.Logger, awsConfig aws.Config, config config.Config, db *sql.DB) (*Factory, error) {
	authClient := newAuthClient(logger, &awsConfig, config)
	domainAuth := newDomainAuth(authClient)
	authUseCases := newAuthUseCases(domainAuth)

	userRepo := newUserRepo(db, logger)
	userService := newUserService(userRepo)
	userUseCases := newUserUseCases(userService, domainAuth, logger)

	return &Factory{
		Domain: Domain{
			Auth: domainAuth,
			User: userService,
		},
		internal: Internal{
			AuthClient: authClient,
			UserRepo:   userRepo,
		},
		UseCases: UseCases{
			Auth: authUseCases,
			User: userUseCases,
		},
	}, nil
}
