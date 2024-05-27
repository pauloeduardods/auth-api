package factory

import (
	"context"
	"monitoring-system/server/config"
	"monitoring-system/server/domain/auth"
	"monitoring-system/server/internal/auth_cognito"
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
	Auth auth.CognitoAuth
}

func newDomainAuth(cognitoAuth auth.CognitoAuth) auth.Auth {
	return auth.New(cognitoAuth)
}

func newInternalAuth(ctx context.Context, logger logger.Logger, awsConfig *aws.Config, config config.Config) auth.CognitoAuth {
	cognitoClient := cognitoidentityprovider.NewFromConfig(*awsConfig)
	jwtVerify := jwt_verify.NewAuth(config.Aws.Region, config.Aws.CognitoUserPoolID, logger)
	jwtVerify.CacheJWK() //TODO: Check when we need to cache the JWK and how to handle the error
	return auth_cognito.NewCognitoAuth(ctx, cognitoClient, config.Aws.CognitoClientId, jwtVerify)
}
func New(ctx context.Context, logger logger.Logger, awsConfig aws.Config, config config.Config) (*Factory, error) {
	internalAuth := newInternalAuth(ctx, logger, &awsConfig, config)
	domainAuth := newDomainAuth(internalAuth)

	return &Factory{
		Domain: Domain{
			Auth: domainAuth,
		},
		internal: Internal{
			Auth: internalAuth,
		},
	}, nil
}
