package factory

import (
	"auth-api/src/config"
	"auth-api/src/internal/domain/admin"
	"auth-api/src/internal/domain/auth"
	"auth-api/src/internal/domain/code"
	"auth-api/src/internal/domain/email"
	"auth-api/src/internal/domain/user"
	"auth-api/src/internal/events"
	eventsIplm "auth-api/src/internal/events"
	events_handlers "auth-api/src/internal/events/handlers"
	admin_repo "auth-api/src/internal/infra/admin/repository"
	admin_service "auth-api/src/internal/infra/admin/service"
	auth_service "auth-api/src/internal/infra/auth/service"
	code_repo "auth-api/src/internal/infra/code/repository"
	code_service "auth-api/src/internal/infra/code/service"
	email_service "auth-api/src/internal/infra/email/service"
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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

type Factory struct {
	UseCases   UseCases
	Repository Repository
	Service    Service
	Event      events.EventDispatcher
}

type Service struct {
	Auth  auth.AuthService
	User  user.UserService
	Admin admin.AdminService
	Code  code.CodeService
	Email email.EmailService
}

type Repository struct {
	User  user.UserRepository
	Admin admin.AdminRepository
	Code  code.CodeRepository
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

func newCodeRepository(awsConfig aws.Config, logger logger.Logger, config config.Config) code.CodeRepository {
	dynamoDBClient := dynamodb.NewFromConfig(awsConfig)
	return code_repo.NewCodeRepositoryDynamoDB(config.Aws.CodesTable, dynamoDBClient, logger)
}

func newEmailService(awsConfig aws.Config, logger logger.Logger) email.EmailService {
	sesClient := ses.NewFromConfig(awsConfig)
	return email_service.NewEmailService(sesClient, logger)
}

func New(ctx context.Context, logger logger.Logger, awsConfig aws.Config, config config.Config, db *sql.DB) (*Factory, error) {
	userRepo := user_repo.NewUserRepository(db, logger)
	adminRepo := admin_repo.NewAdminRepository(db, logger)
	codeRepo := newCodeRepository(awsConfig, logger, config)

	authService := newAuthService(logger, &awsConfig, config)
	userService := user_service.NewUserService(userRepo)
	adminService := admin_service.NewAdminService(adminRepo, logger)
	codeService := code_service.NewCodeServiceImpl(codeRepo, logger)
	emailService := newEmailService(awsConfig, logger)

	dispatcher := eventsIplm.NewEventDispatcher(logger)

	authUseCases := auth_usecases.NewUseCases(authService, adminService, userService, logger, codeService, emailService)
	adminUseCases := admin_usecases.NewUseCases(adminService, authService, logger)
	userUseCases := user_usecases.NewUseCases(userService, authService, logger, dispatcher)

	handlers := events_handlers.NewEventsHandlers(logger, *authUseCases)
	handlers.RegisterHandlers(dispatcher)

	return &Factory{
		Repository: Repository{
			User:  userRepo,
			Admin: adminRepo,
			Code:  codeRepo,
		},
		Service: Service{
			Auth:  authService,
			User:  userService,
			Admin: adminService,
			Code:  codeService,
			Email: emailService,
		},
		UseCases: UseCases{
			Auth:  authUseCases,
			User:  userUseCases,
			Admin: adminUseCases,
		},
		Event: dispatcher,
	}, nil
}
