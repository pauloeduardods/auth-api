package handlers

import (
	"context"
	"monitoring-system/server/src/domain/auth"
	"monitoring-system/server/src/pkg/app_error"
	usecaseUser "monitoring-system/server/src/usecases/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	useCases *usecaseUser.UseCases
}

func NewUserHandler(useCases *usecaseUser.UseCases) *UserHandler {
	return &UserHandler{
		useCases: useCases,
	}
}

type registerInput struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Name     string  `json:"name"`
	Phone    *string `json:"phone"`
}

func (h *UserHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, registerInput{}, func(ctx context.Context, input registerInput) error {
			err := h.useCases.Register.Execute(ctx, usecaseUser.RegisterUserInput{
				Email:    input.Email,
				Password: input.Password,
				Name:     input.Name,
				Phone:    input.Phone,
			})
			return err
		})
	}
}

type updateUserInput struct {
	Name  *string `json:"name"`
	Phone *string `json:"phone"`
}

func (h *UserHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		userId := claims.(*auth.Claims).Id
		if !exists {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			c.Abort()
			return
		}

		processRequestNoOutput(c, updateUserInput{}, func(ctx context.Context, input updateUserInput) error {
			err := h.useCases.Update.Execute(ctx, usecaseUser.UpdateUserInput{
				Id:    userId,
				Name:  input.Name,
				Phone: input.Phone,
			})
			return err
		})
	}
}
