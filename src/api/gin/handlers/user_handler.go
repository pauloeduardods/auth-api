package handlers

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/internal/modules/user-manager/domain/user"
	user_usecases "auth-api/src/internal/modules/user-manager/usecases/user"
	"auth-api/src/pkg/app_error"
	"context"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	useCases *user_usecases.UseCases
}

func NewUserHandler(useCases *user_usecases.UseCases) *UserHandler {
	return &UserHandler{
		useCases: useCases,
	}
}

type registerUserInput struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Name     string  `json:"name"`
	Phone    *string `json:"phone"`
}

func (h *UserHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, registerUserInput{}, func(ctx context.Context, input registerUserInput) error {
			err := h.useCases.Register.Execute(ctx, user_usecases.RegisterUserInput{
				SignUpInput: auth.SignUpInput{
					Username: input.Email,
					Password: input.Password,
					Name:     input.Name,
				},
				CreateUserInput: user.CreateUserInput{
					Phone: input.Phone,
					Name:  input.Name,
					Email: input.Email,
				},
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
		if !exists {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			c.Abort()
			return
		}
		userClaims, ok := claims.(*auth.Claims)
		if !ok {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			c.Abort()
			return
		}
		userId := userClaims.Id

		processRequestNoOutput(c, updateUserInput{}, func(ctx context.Context, input updateUserInput) error {
			userId, err := user.ParseUserID(userId)
			if err != nil {
				return err
			}
			err = h.useCases.Update.Execute(ctx, user_usecases.UpdateUserInput{
				UpdateUserInput: user.UpdateUserInput{
					ID:    userId,
					Name:  input.Name,
					Phone: input.Phone,
				},
			})
			return err
		})
	}
}
