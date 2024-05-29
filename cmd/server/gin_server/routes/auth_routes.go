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
	authGroup.POST("/login", handler.Login())
	authGroup.POST("/refresh", handler.RefreshToken())
	authGroup.POST("/confirm", handler.ConfirmSignUp())

	mfaGroup := authGroup.Group("/mfa")
	mfaGroup.POST("/", handler.AddMfa())
	mfaGroup.POST("/verify", handler.VerifyMfa())
	mfaGroup.POST("/remove", r.authMiddleware.AuthMiddleware(auth.Admin), handler.RemoveMfa())

	authenticatedGroup := authGroup.Group("/")
	authenticatedGroup.Use(r.authMiddleware.AuthMiddleware(auth.Admin, auth.User))
	authenticatedGroup.GET("/", handler.GetUser())

	userGroup := authGroup.Group("/user")
	userGroup.POST("/register", handler.SignUp())

	adminGroup := authGroup.Group("/admin")
	adminGroup.Use(r.authMiddleware.AuthMiddleware(auth.Admin))
	adminGroup.POST("/register", handler.CreateAdmin())

}
