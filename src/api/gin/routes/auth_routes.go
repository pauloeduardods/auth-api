package routes

import (
	"auth-api/src/api/gin/handlers"
	"auth-api/src/api/gin/middleware"
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"time"
)

func (r *routes) configAuthRoutes() {
	handler := handlers.NewAuthHandler(r.factory.UseCases.UserManager.Auth)
	authGroup := r.gin.Group("/auth")
	authGroup.Use(middleware.TimeoutMiddleware(30 * time.Second))

	authGroup.POST("/login", handler.Login())
	authGroup.POST("/logout", handler.Logout())
	authGroup.POST("/refresh", handler.RefreshToken())
	authGroup.POST("/confirm", handler.ConfirmSignUp())
	authGroup.POST("/send-confirmation-code", handler.SendConfirmationCode())
	authGroup.POST("/password/forget", handler.SendForgotPasswordCode())
	authGroup.POST("/password/reset", handler.ResetPassword())
	authGroup.POST("/password/change", handler.ChangePassword())
	authGroup.POST("/password/set", handler.SetPassword())

	mfaGroup := authGroup.Group("/mfa")
	mfaGroup.POST("/", handler.AddMfa())
	mfaGroup.POST("/verify", handler.VerifyMfa())
	mfaGroup.POST("/remove", handler.RemoveMfa())
	mfaGroup.POST("/admin/remove", r.authMiddleware.AuthMiddleware(auth.GroupAdmin), handler.AdminRemoveMfa())
	mfaGroup.POST("/activate", r.authMiddleware.AuthMiddleware(auth.GroupUser), handler.ActivateMfa())

	groupsGroup := authGroup.Group("/groups")
	groupsGroup.POST("/add", r.authMiddleware.AuthMiddleware(auth.GroupAdmin), handler.AddGroup())
	groupsGroup.POST("/remove", r.authMiddleware.AuthMiddleware(auth.GroupAdmin), handler.RemoveGroup())

	authenticatedGroup := authGroup.Group("/")
	authenticatedGroup.Use(r.authMiddleware.AuthMiddleware(auth.GroupAdmin, auth.GroupUser))
	authenticatedGroup.GET("/", handler.GetMe())
}
