package routes

import (
	"monitoring-system/server/cmd/factory"
	"monitoring-system/server/cmd/server/gin_server/middleware"
	"monitoring-system/server/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Routes interface {
	ConfigRoutes()
}

type routes struct {
	gin            *gin.RouterGroup
	factory        *factory.Factory
	validator      validator.Validator
	authMiddleware middleware.AuthMiddleware
}

func NewRoutes(g *gin.RouterGroup, factory *factory.Factory, v validator.Validator, authMiddleware middleware.AuthMiddleware) Routes {
	return &routes{
		gin:            g,
		factory:        factory,
		validator:      v,
		authMiddleware: authMiddleware,
	}
}

func (r *routes) ConfigRoutes() {
	r.configAuthRoutes()
}
