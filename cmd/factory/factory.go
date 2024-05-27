package factory

import (
	"context"
	"monitoring-system/server/config"
	"monitoring-system/server/domain/auth"
	"monitoring-system/server/internal/auth_cognito"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type Factory struct {
	Domain   Domain
	Internal Internal
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

func newInternalAuth(ctx context.Context, awsConfig *aws.Config, config config.Config) auth.CognitoAuth {
	cognitoClient := cognitoidentityprovider.NewFromConfig(*awsConfig)
	return auth_cognito.NewCognitoAuth(ctx, cognitoClient, config.Aws.CognitoClientId)
}
func New(ctx context.Context, awsConfig aws.Config, config config.Config) (*Factory, error) {
	internalAuth := newInternalAuth(ctx, &awsConfig, config)
	domainAuth := newDomainAuth(internalAuth)

	return &Factory{
		Domain: Domain{
			Auth: domainAuth,
		},
		Internal: Internal{
			Auth: internalAuth,
		},
	}, nil
}
