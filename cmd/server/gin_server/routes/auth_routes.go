package routes

import (
	"monitoring-system/server/cmd/server/gin_server/handlers"
)

type AuthRoutes struct {
}

func (r *routes) configAuthRoutes() {
	handler := handlers.NewAuthHandler(r.factory.Domain.Auth, r.validator)
	authGroup := r.gin.Group("/auth")

	authGroup.POST("/login", handler.Login())
	authGroup.POST("/register", handler.SignUp())
	authGroup.POST("/confirm", handler.ConfirmSignUp())
	authGroup.GET("/user", r.authMiddleware.AuthMiddleware(), handler.GetUser())
}
