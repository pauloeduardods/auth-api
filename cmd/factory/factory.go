package factory

import (
	"context"
	"monitoring-system/server/config"
	authDomain "monitoring-system/server/domain/auth"
	auth_client "monitoring-system/server/internal/auth/client"
	"monitoring-system/server/pkg/jwt_verify"
	"monitoring-system/server/pkg/logger"
	authUseCases "monitoring-system/server/usecases/auth"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type Factory struct {
	Domain   Domain
	internal Internal
	UseCases UseCases
}

type Domain struct {
	Auth authDomain.Auth
}

type Internal struct {
	AuthClient authDomain.AuthClient
}

type UseCases struct {
	Auth *authUseCases.UseCases
}

func newDomainAuth(authClient authDomain.AuthClient) authDomain.Auth {
	return authDomain.NewAuthService(authClient)
}

func newAuthClient(logger logger.Logger, awsConfig *aws.Config, config config.Config) authDomain.AuthClient {
	cognitoClient := cognitoidentityprovider.NewFromConfig(*awsConfig)
	jwtVerify := jwt_verify.NewAuth(config.Aws.Region, config.Aws.CognitoUserPoolID, logger)
	jwtVerify.CacheJWK() //TODO: Check when we need to cache the JWK and how to handle the error
	return auth_client.NewAuthClient(cognitoClient, config.Aws.CognitoClientId, jwtVerify, config.Aws.CognitoUserPoolID, logger)
}

func newAuthUseCases(auth authDomain.Auth) *authUseCases.UseCases {
	return authUseCases.NewUseCases(auth)
}

func New(ctx context.Context, logger logger.Logger, awsConfig aws.Config, config config.Config) (*Factory, error) {
	authClient := newAuthClient(logger, &awsConfig, config)
	domainAuth := newDomainAuth(authClient)
	authUseCases := newAuthUseCases(domainAuth)

	return &Factory{
		Domain: Domain{
			Auth: domainAuth,
		},
		internal: Internal{
			AuthClient: authClient,
		},
		UseCases: UseCases{
			Auth: authUseCases,
		},
	}, nil
}
