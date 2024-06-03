package routes

import (
	"auth-api/src/api/gin/handlers"
	"auth-api/src/api/gin/middleware"
	"auth-api/src/internal/domain/auth"
	"time"
)

func (r *routes) configAuthRoutes() {
	handler := handlers.NewAuthHandler(r.factory.UseCases.Auth)
	authGroup := r.gin.Group("/auth")
	authGroup.Use(middleware.TimeoutMiddleware(30 * time.Second))

	authGroup.POST("/login", handler.Login())
	authGroup.POST("/logout", handler.Logout())
	authGroup.POST("/refresh", handler.RefreshToken())
	authGroup.POST("/confirm", handler.ConfirmSignUp())
	authGroup.POST("/set-password", handler.SetPassword())

	mfaGroup := authGroup.Group("/mfa")
	mfaGroup.POST("/", handler.AddMfa())
	mfaGroup.POST("/verify", handler.VerifyMfa())
	mfaGroup.POST("/remove", handler.RemoveMfa())
	mfaGroup.POST("/admin/remove", r.authMiddleware.AuthMiddleware(auth.Admin), handler.AdminRemoveMfa())
	mfaGroup.POST("/activate", r.authMiddleware.AuthMiddleware(auth.User), handler.ActivateMfa())

	groupsGroup := authGroup.Group("/groups")
	groupsGroup.POST("/add", r.authMiddleware.AuthMiddleware(auth.Admin), handler.AddGroup())
	groupsGroup.POST("/remove", r.authMiddleware.AuthMiddleware(auth.Admin), handler.RemoveGroup())

	authenticatedGroup := authGroup.Group("/")
	authenticatedGroup.Use(r.authMiddleware.AuthMiddleware(auth.Admin, auth.User))
	authenticatedGroup.GET("/", handler.GetMe())
}
