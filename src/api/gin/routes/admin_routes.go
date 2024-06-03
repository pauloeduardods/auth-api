package routes

import (
	"auth-api/src/api/gin/handlers"
	"auth-api/src/api/gin/middleware"
	"auth-api/src/internal/domain/auth"
	"time"
)

func (r *routes) configAdminRoutes() {
	handler := handlers.NewAdminHandler(r.factory.UseCases.Admin)
	adminGroup := r.gin.Group("/admin")
	adminGroup.Use(middleware.TimeoutMiddleware(30 * time.Second))

	adminGroup.PATCH("/", r.authMiddleware.AuthMiddleware(auth.Admin), handler.Update())
	adminGroup.POST("/register", r.authMiddleware.AuthMiddleware(auth.Admin), handler.Register())

}
