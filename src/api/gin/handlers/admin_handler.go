package handlers

import (
	"auth-api/src/internal/domain/admin"
	"auth-api/src/internal/domain/auth"
	admin_usecases "auth-api/src/internal/usecases/admin"
	"auth-api/src/pkg/app_error"
	"context"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	useCases *admin_usecases.UseCases
}

func NewAdminHandler(useCases *admin_usecases.UseCases) *AdminHandler {
	return &AdminHandler{
		useCases: useCases,
	}
}

type registerAdminInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h *AdminHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, registerAdminInput{}, func(ctx context.Context, input registerAdminInput) error {
			err := h.useCases.Register.Execute(ctx, admin_usecases.RegisterAdminInput{
				SignupAdmin: auth.CreateAdminInput{
					Username: input.Email,
					Password: input.Password,
					Name:     input.Name,
				},
				CreateAdminInput: admin.CreateAdminInput{
					Name:  input.Name,
					Email: input.Email,
				},
			})
			return err
		})
	}
}

type updateAdminInput struct {
	Name *string `json:"name"`
}

func (h *AdminHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		adminId := claims.(*auth.Claims).Id
		if !exists {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			c.Abort()
			return
		}

		processRequestNoOutput(c, updateAdminInput{}, func(ctx context.Context, input updateAdminInput) error {
			adminId, err := admin.ParseAdminID(adminId)
			if err != nil {
				return err
			}
			err = h.useCases.Update.Execute(ctx, admin_usecases.UpdateAdminInput{
				UpdateAdminInput: admin.UpdateAdminInput{
					ID:   adminId,
					Name: input.Name,
				},
			})
			return err
		})
	}
}
