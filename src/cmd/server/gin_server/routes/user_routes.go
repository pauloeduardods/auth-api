package routes

import (
	"monitoring-system/server/src/cmd/server/gin_server/handlers"
	"monitoring-system/server/src/cmd/server/gin_server/middleware"
	"monitoring-system/server/src/domain/auth"
	"time"
)

func (r *routes) configUserRoutes() {
	handler := handlers.NewUserHandler(r.factory.UseCases.User)
	userGroup := r.gin.Group("/user")
	userGroup.Use(middleware.TimeoutMiddleware(30 * time.Second))

	userGroup.PATCH("/", r.authMiddleware.AuthMiddleware(auth.User), handler.Update())
	userGroup.POST("/register", handler.Register())

}
