package factory

import (
	"context"
	"monitoring-system/server/config"
	"monitoring-system/server/domain/auth"
	auth_client "monitoring-system/server/internal/auth/client"
	"monitoring-system/server/pkg/jwt_verify"
	"monitoring-system/server/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type Factory struct {
	Domain   Domain
	internal Internal
}

type Domain struct {
	Auth auth.Auth
}

type Internal struct {
	AuthClient auth.AuthClient
}

func newDomainAuth(authClient auth.AuthClient) auth.Auth {
	return auth.NewAuthService(authClient)
}

func newAuthClient(ctx context.Context, logger logger.Logger, awsConfig *aws.Config, config config.Config) auth.AuthClient {
	cognitoClient := cognitoidentityprovider.NewFromConfig(*awsConfig)
	jwtVerify := jwt_verify.NewAuth(config.Aws.Region, config.Aws.CognitoUserPoolID, logger)
	jwtVerify.CacheJWK() //TODO: Check when we need to cache the JWK and how to handle the error
	return auth_client.NewAuthClient(ctx, cognitoClient, config.Aws.CognitoClientId, jwtVerify, config.Aws.CognitoUserPoolID, logger)
}
func New(ctx context.Context, logger logger.Logger, awsConfig aws.Config, config config.Config) (*Factory, error) {
	authClient := newAuthClient(ctx, logger, &awsConfig, config)
	domainAuth := newDomainAuth(authClient)

	return &Factory{
		Domain: Domain{
			Auth: domainAuth,
		},
		internal: Internal{
			AuthClient: authClient,
		},
	}, nil
}
