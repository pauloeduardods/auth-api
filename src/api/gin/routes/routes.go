package routes

import (
	"auth-api/src/api/gin/middleware"
	"auth-api/src/factory"

	"github.com/gin-gonic/gin"
)

type Routes interface {
	ConfigRoutes()
}

type routes struct {
	gin            *gin.RouterGroup
	factory        *factory.Factory
	authMiddleware middleware.AuthMiddleware
}

func NewRoutes(g *gin.RouterGroup, factory *factory.Factory, authMiddleware middleware.AuthMiddleware) Routes {
	return &routes{
		gin:            g,
		factory:        factory,
		authMiddleware: authMiddleware,
	}
}

func (r *routes) ConfigRoutes() {
	r.configAuthRoutes()
	r.configUserRoutes()
}
