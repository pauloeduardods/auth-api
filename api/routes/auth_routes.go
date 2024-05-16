package routes

import (
	handler "monitoring-system/server/api/handlers"
	"monitoring-system/server/api/middleware"

	"github.com/gin-gonic/gin"
)

func ConfigAuthRoutes(g *gin.Engine, m middleware.AuthMiddleware, h handler.AuthHandler) {
	authGroup := g.Group("/api/v1/auth")

	authGroup.POST("/login", h.Login())
	authGroup.POST("/register", h.Register())
	authGroup.POST("/confirm", h.ConfirmSignUp())
	authGroup.GET("/info", m.AuthMiddleware(), h.GetUser())
}
