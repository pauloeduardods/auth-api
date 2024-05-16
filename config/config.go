package config

import (
	"context"
	"fmt"
	"monitoring-system/server/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/logging"
	"github.com/maragudk/env"
)

type AppConfig struct {
	Host              string
	Port              int
	CognitoClientId   string
	Region            string
	CognitoUserPoolID string
	AppEnv            string
}

func LoadConfig() (*AppConfig, error) {
	_ = env.Load(".env")

	appEnv := env.GetStringOrDefault("APP_ENV", "development")
	host := env.GetStringOrDefault("HOST", "localhost")
	port := env.GetIntOrDefault("PORT", 4000)
	cognitoClientId := env.GetStringOrDefault("COGNITO_CLIENT_ID", "")
	cognitoUserPoolID := env.GetStringOrDefault("COGNITO_USER_POOL_ID", "")
	region := env.GetStringOrDefault("REGION", "")

	if host == "" {
		return nil, fmt.Errorf("host is not defined")
	}
	if cognitoClientId == "" {
		return nil, fmt.Errorf("cognitoClientId is not defined")
	}
	if cognitoUserPoolID == "" {
		return nil, fmt.Errorf("CognitoUserPoolID is not defined")
	}
	if region == "" {
		return nil, fmt.Errorf("region is not defined")
	}

	config := &AppConfig{
		Host:              host,
		Port:              port,
		CognitoClientId:   cognitoClientId,
		Region:            region,
		CognitoUserPoolID: cognitoUserPoolID,
		AppEnv:            appEnv,
	}
	return config, nil
}

func NewAWSConfig(ctx context.Context, env *AppConfig, logger logger.Logger) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(env.Region),
		config.WithLogger(createAWSLogAdapter(logger)),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading aws configuration: %v", err)
	}

	return &cfg, nil
}

func createAWSLogAdapter(log logger.Logger) logging.LoggerFunc {
	return func(classification logging.Classification, format string, v ...interface{}) {
		switch classification {
		case logging.Debug:
			log.Debug(format, v...)
		case logging.Warn:
			log.Warning(format, v...)
		default:
			log.Info(format, v...)
		}
	}
}
