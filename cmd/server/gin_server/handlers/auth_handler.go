package handlers

import (
	"io"
	domainAuth "monitoring-system/server/domain/auth"
	"monitoring-system/server/pkg/app_error"
	usecaseAuth "monitoring-system/server/usecases/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	useCases *usecaseAuth.UseCases
}

func NewAuthHandler(useCases *usecaseAuth.UseCases) *AuthHandler {
	return &AuthHandler{
		useCases: useCases,
	}
}

type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input loginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			if err == io.EOF {
				c.Error(app_error.NewApiError(400, "Invalid request"))
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		output, err := h.useCases.Login.Execute(c.Request.Context(), usecaseAuth.LoginInput{
			LoginInput: domainAuth.LoginInput{
				Username: input.Email,
				Password: input.Password,
			},
		})
		if err != nil {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

// type signUpInput struct {
// 	Email    string `json:"email" binding:"required" validate:"email"`
// 	Password string `json:"password" binding:"required" validate:"min=8"`
// 	Name     string `json:"name" binding:"required" validate:"min=3,max=50"`
// }

// func (h *AuthHandler) SignUp() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var input signUpInput
// 		if err := c.ShouldBindJSON(&input); err != nil {
// if err == io.EOF {
// 			c.Error(app_error.NewApiError(400, "Invalid request"))
// 			return
// 		}
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		output, err := h.useCases.SignUp.Execute(c.Request.Context(), usecaseAuth.SignUpInput{
// 			SignUpInput: domainAuth.SignUpInput{
// 				Username: input.Email,
// 				Password: input.Password,
// 				Name:     input.Name,
// 			},
// 		})
// 		if err != nil {
// 			c.Error(app_error.NewApiError(400, "Failed to sign up"))
// 			return
// 		}
// 		c.JSON(http.StatusOK, output)
// 	}
// }

type confirmSignUpInput struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *AuthHandler) ConfirmSignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input confirmSignUpInput
		if err := c.ShouldBindJSON(&input); err != nil {
			if err == io.EOF {
				c.Error(app_error.NewApiError(400, "Invalid request"))
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := h.useCases.ConfirmSignUp.Execute(c.Request.Context(), usecaseAuth.ConfirmSignUpInput{
			ConfirmSignUpInput: domainAuth.ConfirmSignUpInput{
				Username: input.Email,
				Code:     input.Code,
			},
		})
		if err != nil {
			c.Error(app_error.NewApiError(400, "Failed to confirm sign up"))
			return
		}
		c.JSON(http.StatusNoContent, gin.H{})
	}
}

type getMeInput struct {
	AccessToken string `form:"accessToken"`
}

func (h *AuthHandler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input getMeInput
		if err := c.ShouldBindQuery(&input); err != nil {
			if err == io.EOF {
				c.Error(app_error.NewApiError(400, "Invalid request"))
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		output, err := h.useCases.GetMe.Execute(c.Request.Context(), usecaseAuth.GetMeInput{
			GetMeInput: domainAuth.GetMeInput{
				AccessToken: input.AccessToken,
			},
		})
		if err != nil {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

type refreshTokenInput struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *AuthHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input refreshTokenInput
		if err := c.ShouldBindJSON(&input); err != nil {
			if err == io.EOF {
				c.Error(app_error.NewApiError(400, "Invalid request"))
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		output, err := h.useCases.RefreshToken.Execute(c.Request.Context(), usecaseAuth.RefreshTokenInput{
			RefreshTokenInput: domainAuth.RefreshTokenInput{
				RefreshToken: input.RefreshToken,
			},
		})
		if err != nil {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

// type createAdminInput struct {
// 	Email    string `json:"email" binding:"required" validate:"email"`
// 	Password string `json:"password" binding:"required" validate:"min=8"`
// 	Name     string `json:"name" binding:"required" validate:"min=3,max=50"`
// }

// func (h *AuthHandler) CreateAdmin() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var input createAdminInput
// 		if err := c.ShouldBindJSON(&input); err != nil {
// if err == io.EOF {
// 			c.Error(app_error.NewApiError(400, "Invalid request"))
// 			return
// 		}
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		output, err := h.useCases.CreateAdmin.Execute(c.Request.Context(), usecaseAuth.CreateAdminInput{
// 			CreateAdminInput: domainAuth.CreateAdminInput{
// 				Username: input.Email,
// 				Password: input.Password,
// 				Name:     input.Name,
// 			},
// 		})
// 		if err != nil {
// 			c.Error(app_error.NewApiError(400, "Failed to create admin"))
// 			return
// 		}
// 		c.JSON(http.StatusOK, output)
// 	}
// }

type addMfaInput struct {
	AccessToken string `json:"accessToken"`
}

func (h *AuthHandler) AddMfa() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input addMfaInput
		if err := c.ShouldBindJSON(&input); err != nil {
			if err == io.EOF {
				c.Error(app_error.NewApiError(400, "Invalid request"))
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		output, err := h.useCases.AddMFA.Execute(c.Request.Context(), usecaseAuth.AddMFAInput{
			AddMFAInput: domainAuth.AddMFAInput{
				AccessToken: input.AccessToken,
			},
		})
		if err != nil {
			c.Error(app_error.NewApiError(400, "Failed to add MFA"))
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

type verifyMfaInput struct {
	Email   string `json:"email"`
	Code    string `json:"code"`
	Session string `json:"session"`
}

func (h *AuthHandler) VerifyMfa() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input verifyMfaInput
		if err := c.ShouldBindJSON(&input); err != nil {
			if err == io.EOF {
				c.Error(app_error.NewApiError(400, "Invalid request"))
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		output, err := h.useCases.VerifyMFA.Execute(c.Request.Context(), usecaseAuth.VerifyMFAInput{
			VerifyMFAInput: domainAuth.VerifyMFAInput{
				Code:     input.Code,
				Username: input.Email,
				Session:  input.Session,
			},
		})
		if err != nil {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

type removeMfaInput struct {
	Username string `json:"username" binding:"required"`
}

func (h *AuthHandler) RemoveMfa() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input removeMfaInput
		if err := c.ShouldBindJSON(&input); err != nil {
			if err == io.EOF {
				c.Error(app_error.NewApiError(400, "Invalid request"))
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := h.useCases.RemoveMFA.Execute(c.Request.Context(), usecaseAuth.RemoveMFAInput{
			RemoveMFAInput: domainAuth.RemoveMFAInput{
				Username: input.Username,
			},
		})
		if err != nil {
			c.Error(app_error.NewApiError(400, "Failed to remove MFA"))
			return
		}
		c.JSON(http.StatusNoContent, gin.H{})
	}
}
