package gin_server

import (
	"context"
	"monitoring-system/server/cmd/factory"
	"monitoring-system/server/cmd/server/gin_server/middleware"
	"monitoring-system/server/cmd/server/gin_server/routes"
	"monitoring-system/server/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	log     logger.Logger
	Gin     *gin.Engine
	ctx     context.Context
	factory *factory.Factory
}

func New(ctx context.Context, logger logger.Logger, factory *factory.Factory) *Gin {
	gin := gin.Default()
	return &Gin{
		log:     logger,
		Gin:     gin,
		ctx:     ctx,
		factory: factory,
	}
}

func (s *Gin) SetupMiddlewares() {
	cors := middleware.NewCors("*", "GET, POST, PUT, DELETE, OPTIONS", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Auth-Token, X-Requested-With", false)
	s.Gin.Use(cors.CorsMiddleware())
	s.Gin.Use(gin.CustomRecovery(middleware.RecoveryHandler(s.log)))
	s.Gin.Use(gin.Logger())
	s.Gin.Use(middleware.ErrorHandler(s.log))
}

func (s *Gin) SetupApi() error {
	//Api Routes
	s.Gin.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	apiRoutes := s.Gin.Group("/api/v1")

	// Middlewares
	authMiddleware := middleware.NewAuthMiddleware(s.factory.Domain.Auth)

	//Static files
	s.Gin.StaticFS("/web", http.Dir("static"))

	//Routes
	routes.NewRoutes(apiRoutes, s.factory, authMiddleware).ConfigRoutes()
	return nil
}
