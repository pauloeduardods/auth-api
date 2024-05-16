package server

import (
	"context"
	"monitoring-system/server/config"
	"monitoring-system/server/pkg/jwtToken"
	"monitoring-system/server/pkg/logger"
	"monitoring-system/server/pkg/validator"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
)

type Server struct {
	log           logger.Logger
	appConfig     *config.AppConfig
	gin           *gin.Engine
	cognitoClient *cognito.Client
	server        *http.Server
	jwtToken      *jwtToken.JwtToken
	validator     validator.Validator
	ctx           context.Context
}

func New(ctx context.Context, awsConfig *aws.Config, appConfig *config.AppConfig, logger logger.Logger) *Server {
	gin := gin.Default()

	cognito := cognito.New(cognito.Options{
		Credentials: awsConfig.Credentials,
		Region:      awsConfig.Region,
	})

	jwtService := jwtToken.NewAuth(appConfig.Region, appConfig.CognitoUserPoolID, logger)

	return &Server{
		appConfig:     appConfig,
		gin:           gin,
		cognitoClient: cognito,
		jwtToken:      jwtService,
		log:           logger,
		validator:     validator.NewValidatorImpl(),
		ctx:           ctx,
	}
}

func (s *Server) Start() error {
	s.log.Info("Starting server %s:%d", s.appConfig.Host, s.appConfig.Port)
	s.SetupCors()
	s.SetupMiddlewares()
	s.SetupApi()

	go func() {
		<-s.ctx.Done()
		s.log.Info("Shutdown Server ...")

		if err := s.server.Shutdown(s.ctx); err != nil {
			s.log.Error("Server Shutdown: %v", err)
		}
		s.log.Info("Server exiting")
	}()

	s.server = &http.Server{
		Addr:    s.appConfig.Host + ":" + strconv.Itoa(s.appConfig.Port),
		Handler: s.gin,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	if s.server != nil {
		s.log.Info("Stopping server")
		if err := s.server.Shutdown(s.ctx); err != nil {
			return err
		}
	}
	return nil
}
