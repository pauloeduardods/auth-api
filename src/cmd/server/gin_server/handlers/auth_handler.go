package handlers

import (
	"context"
	domainAuth "monitoring-system/server/src/domain/auth"
	usecaseAuth "monitoring-system/server/src/usecases/auth"

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
		processRequest(c, loginInput{}, func(ctx context.Context, input loginInput) (*domainAuth.LoginOutput, error) {
			return h.useCases.Login.Execute(ctx, usecaseAuth.LoginInput{
				LoginInput: domainAuth.LoginInput{
					Username: input.Email,
					Password: input.Password,
				},
			})
		})
	}
}

type confirmSignUpInput struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *AuthHandler) ConfirmSignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, confirmSignUpInput{}, func(ctx context.Context, input confirmSignUpInput) error {
			_, err := h.useCases.ConfirmSignUp.Execute(ctx, usecaseAuth.ConfirmSignUpInput{
				ConfirmSignUpInput: domainAuth.ConfirmSignUpInput{
					Username: input.Email,
					Code:     input.Code,
				},
			})
			return err
		})
	}
}

type getMeInput struct {
	AccessToken string `form:"accessToken"`
}

func (h *AuthHandler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestQuery(c, getMeInput{}, func(ctx context.Context, input getMeInput) (*domainAuth.GetMeOutput, error) {
			return h.useCases.GetMe.Execute(ctx, usecaseAuth.GetMeInput{
				GetMeInput: domainAuth.GetMeInput{
					AccessToken: input.AccessToken,
				},
			})
		})
	}
}

type refreshTokenInput struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *AuthHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequest(c, refreshTokenInput{}, func(ctx context.Context, input refreshTokenInput) (*domainAuth.RefreshTokenOutput, error) {
			return h.useCases.RefreshToken.Execute(ctx, usecaseAuth.RefreshTokenInput{
				RefreshTokenInput: domainAuth.RefreshTokenInput{
					RefreshToken: input.RefreshToken,
				},
			})
		})
	}
}

type addMfaInput struct {
	AccessToken string `json:"accessToken"`
}

func (h *AuthHandler) AddMfa() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequest(c, addMfaInput{}, func(ctx context.Context, input addMfaInput) (*domainAuth.AddMFAOutput, error) {
			return h.useCases.AddMFA.Execute(ctx, usecaseAuth.AddMFAInput{
				AddMFAInput: domainAuth.AddMFAInput{
					AccessToken: input.AccessToken,
				},
			})
		})
	}
}

type verifyMfaInput struct {
	Email   string `json:"email"`
	Code    string `json:"code"`
	Session string `json:"session"`
}

func (h *AuthHandler) VerifyMfa() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequest(c, verifyMfaInput{}, func(ctx context.Context, input verifyMfaInput) (*domainAuth.LoginOutput, error) {
			return h.useCases.VerifyMFA.Execute(ctx, usecaseAuth.VerifyMFAInput{
				VerifyMFAInput: domainAuth.VerifyMFAInput{
					Code:     input.Code,
					Username: input.Email,
					Session:  input.Session,
				},
			})
		})
	}
}

type adminRemoveMfaInput struct {
	Username string `json:"username"`
}

func (h *AuthHandler) AdminRemoveMfa() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, adminRemoveMfaInput{}, func(ctx context.Context, input adminRemoveMfaInput) error {
			return h.useCases.AdminRemoveMFA.Execute(ctx, usecaseAuth.AdminRemoveMFAInput{
				AdminRemoveMFAInput: domainAuth.AdminRemoveMFAInput{
					Username: input.Username,
				},
			})
		})
	}
}

type removeMfaInput struct {
	AccessToken string `json:"accessToken"`
}

func (h *AuthHandler) RemoveMfa() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, removeMfaInput{}, func(ctx context.Context, input removeMfaInput) error {
			return h.useCases.RemoveMFA.Execute(ctx, usecaseAuth.RemoveMFAInput{
				RemoveMFAInput: domainAuth.RemoveMFAInput{
					AccessToken: input.AccessToken,
				},
			})
		})
	}
}

type activateMfaInput struct {
	AccessToken string `json:"accessToken"`
	Code        string `json:"code"`
}

func (h *AuthHandler) ActivateMfa() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, activateMfaInput{}, func(ctx context.Context, input activateMfaInput) error {
			err := h.useCases.ActivateMFA.Execute(ctx, usecaseAuth.ActivateMFAInput{
				ActivateMFAInput: domainAuth.ActivateMFAInput{
					AccessToken: input.AccessToken,
					Code:        input.Code,
				},
			})
			return err
		})
	}
}

type logoutInput struct {
	AccessToken string `json:"accessToken"`
}

func (h *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, logoutInput{}, func(ctx context.Context, input logoutInput) error {
			err := h.useCases.Logout.Execute(ctx, usecaseAuth.LogoutInput{
				LogoutInput: domainAuth.LogoutInput{
					AccessToken: input.AccessToken,
				},
			})
			return err
		})
	}
}
