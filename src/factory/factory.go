package factory

import (
	"auth-api/src/config"
	"auth-api/src/internal/events"
	eventsIplm "auth-api/src/internal/events"
	events_handlers "auth-api/src/internal/events/handlers"
	"auth-api/src/internal/modules/user-manager/domain/admin"
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/internal/modules/user-manager/domain/user"
	admin_infra "auth-api/src/internal/modules/user-manager/infra/admin"
	auth_infra "auth-api/src/internal/modules/user-manager/infra/auth"
	user_infra "auth-api/src/internal/modules/user-manager/infra/user"
	admin_usecases "auth-api/src/internal/modules/user-manager/usecases/admin"
	auth_usecases "auth-api/src/internal/modules/user-manager/usecases/auth"
	user_usecases "auth-api/src/internal/modules/user-manager/usecases/user"
	"auth-api/src/internal/shared/code/domain/code"
	code_infra "auth-api/src/internal/shared/code/infra/code"
	"auth-api/src/internal/shared/notification/domain/email"
	email_infra "auth-api/src/internal/shared/notification/infra/email"
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
	Code        code.CodeService
	Email       email.EmailService
	UserManager UserManagerService
}

type Repository struct {
	UserManager UserManagerRepo
	Code        code.CodeRepository
}

type UserManagerService struct {
	Auth  auth.AuthService
	User  user.UserService
	Admin admin.AdminService
}

type UserManagerRepo struct {
	User  user.UserRepository
	Admin admin.AdminRepository
}

type UserManagerUseCases struct {
	Auth  *auth_usecases.UseCases
	User  *user_usecases.UseCases
	Admin *admin_usecases.UseCases
}

type UseCases struct {
	UserManager UserManagerUseCases
}

func newAuthService(logger logger.Logger, awsConfig *aws.Config, config config.Config, email email.EmailService, codeService code.CodeService) auth.AuthService {
	cognitoClient := cognitoidentityprovider.NewFromConfig(*awsConfig)
	jwtVerify := jwt_verify.NewAuth(config.Aws.Region, config.Aws.CognitoUserPoolID, logger)
	jwtVerify.CacheJWK() //TODO: Check when we need to cache the JWK and how to handle the error
	return auth_infra.NewAuthService(cognitoClient, config.Aws.CognitoClientId, jwtVerify, config.Aws.CognitoUserPoolID, logger, email, codeService)
}

func newCodeRepository(awsConfig aws.Config, logger logger.Logger, config config.Config) code.CodeRepository {
	dynamoDBClient := dynamodb.NewFromConfig(awsConfig)
	return code_infra.NewCodeRepositoryDynamoDB(config.Aws.CodesTable, dynamoDBClient, logger)
}

func newEmailService(awsConfig aws.Config, logger logger.Logger) email.EmailService {
	sesClient := ses.NewFromConfig(awsConfig)
	return email_infra.NewEmailService(sesClient, logger)
}

func New(ctx context.Context, logger logger.Logger, awsConfig aws.Config, config config.Config, db *sql.DB) (*Factory, error) {
	userRepo := user_infra.NewUserRepository(db, logger)
	adminRepo := admin_infra.NewAdminRepository(db, logger)
	codeRepo := newCodeRepository(awsConfig, logger, config)

	codeService := code_infra.NewCodeServiceImpl(codeRepo, logger)
	emailService := newEmailService(awsConfig, logger)

	authService := newAuthService(logger, &awsConfig, config, emailService, codeService)
	userService := user_infra.NewUserService(userRepo)
	adminService := admin_infra.NewAdminService(adminRepo, logger)

	dispatcher := eventsIplm.NewEventDispatcher(logger)

	authUseCases := auth_usecases.NewUseCases(authService, adminService, userService, logger)
	adminUseCases := admin_usecases.NewUseCases(adminService, authService, logger)
	userUseCases := user_usecases.NewUseCases(userService, authService, logger, dispatcher)

	handlers := events_handlers.NewEventsHandlers(logger, *authUseCases)
	handlers.RegisterHandlers(dispatcher)

	return &Factory{
		Repository: Repository{
			UserManager: UserManagerRepo{
				User:  userRepo,
				Admin: adminRepo,
			},
			Code: codeRepo,
		},
		Service: Service{
			UserManager: UserManagerService{
				Auth:  authService,
				User:  userService,
				Admin: adminService,
			},
			Code:  codeService,
			Email: emailService,
		},
		UseCases: UseCases{
			UserManager: UserManagerUseCases{
				Auth:  authUseCases,
				User:  userUseCases,
				Admin: adminUseCases,
			},
		},
		Event: dispatcher,
	}, nil
}
