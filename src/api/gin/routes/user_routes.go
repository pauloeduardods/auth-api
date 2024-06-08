package routes

import (
	"auth-api/src/api/gin/handlers"
	"auth-api/src/api/gin/middleware"
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"time"
)

func (r *routes) configUserRoutes() {
	handler := handlers.NewUserHandler(r.factory.UseCases.UserManager.User)
	userGroup := r.gin.Group("/user")
	userGroup.Use(middleware.TimeoutMiddleware(30 * time.Second))

	userGroup.PATCH("/", r.authMiddleware.AuthMiddleware(auth.GroupUser), handler.Update())
	userGroup.POST("/register", handler.Register())

}
