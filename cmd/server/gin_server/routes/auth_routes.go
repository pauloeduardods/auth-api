package routes

import (
	"monitoring-system/server/cmd/server/gin_server/handlers"
	"monitoring-system/server/domain/auth"
)

type AuthRoutes struct {
}

func (r *routes) configAuthRoutes() {
	handler := handlers.NewAuthHandler(r.factory.Domain.Auth, r.validator)
	authGroup := r.gin.Group("/auth")

	userGroup := authGroup.Group("/user")
	userGroup.GET("/", r.authMiddleware.AuthMiddleware(auth.User), handler.GetUser())
	userGroup.POST("/register", handler.SignUp(auth.User))

	adminGroup := authGroup.Group("/admin")
	adminGroup.Use(r.authMiddleware.AuthMiddleware(auth.Admin))
	adminGroup.POST("/register", handler.SignUp(auth.Admin))

	authGroup.POST("/login", handler.Login())
	authGroup.POST("/confirm", handler.ConfirmSignUp())
}
