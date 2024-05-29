package routes

import (
	"monitoring-system/server/src/cmd/server/gin_server/handlers"
	"monitoring-system/server/src/cmd/server/gin_server/middleware"
	"monitoring-system/server/src/domain/auth"
	"time"
)

func (r *routes) configAuthRoutes() {
	handler := handlers.NewAuthHandler(r.factory.UseCases.Auth)
	authGroup := r.gin.Group("/auth")
	authGroup.Use(middleware.TimeoutMiddleware(30 * time.Second))

	authGroup.POST("/login", handler.Login())
	authGroup.POST("/refresh", handler.RefreshToken())
	authGroup.POST("/confirm", handler.ConfirmSignUp())

	mfaGroup := authGroup.Group("/mfa")
	mfaGroup.POST("/", handler.AddMfa())
	mfaGroup.POST("/verify", handler.VerifyMfa())
	mfaGroup.POST("/remove", r.authMiddleware.AuthMiddleware(auth.Admin), handler.RemoveMfa())

	authenticatedGroup := authGroup.Group("/")
	authenticatedGroup.Use(r.authMiddleware.AuthMiddleware(auth.Admin, auth.User))
	authenticatedGroup.GET("/", handler.GetMe())

	// userGroup := authGroup.Group("/user")
	// userGroup.POST("/register", handler.SignUp())

	// adminGroup := authGroup.Group("/admin")
	// adminGroup.Use(r.authMiddleware.AuthMiddleware(auth.Admin))
	// adminGroup.POST("/register", handler.CreateAdmin())

}
