package server

import (
	"context"
	"monitoring-system/server/src/api/gin"
	"monitoring-system/server/src/cmd/factory"
	"monitoring-system/server/src/config"
	"monitoring-system/server/src/pkg/logger"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type Server struct {
	log    logger.Logger
	config *config.Config
	gin    *gin.Gin
	server *http.Server
	ctx    context.Context
}

func New(ctx context.Context, awsConfig *aws.Config, config *config.Config, logger logger.Logger, factory *factory.Factory) *Server {
	gin := gin.New(logger, factory)

	return &Server{
		config: config,
		gin:    gin,
		log:    logger,
		ctx:    ctx,
	}
}

func (s *Server) Start() error {
	s.log.Info("Starting server %s:%d", s.config.Api.Host, s.config.Api.Port)

	s.gin.SetupMiddlewares()
	s.gin.SetupApi()

	go func() {
		<-s.ctx.Done()
		s.log.Info("Shutdown Server ...")

		if err := s.server.Shutdown(s.ctx); err != nil {
			s.log.Error("Server Shutdown: %v", err)
		}
		s.log.Info("Server exiting")
	}()

	s.server = &http.Server{
		Addr:    s.config.Api.Host + ":" + strconv.Itoa(s.config.Api.Port),
		Handler: s.gin.Gin,
	}

	err := s.server.ListenAndServe()
	if err != nil {
		errType := err.Error()
		if errType == "http: Server closed" {
			s.log.Info("Server closed")
			return nil
		}

		s.log.Error("Error starting server: %v", err)
		return err
	}
	return nil
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
