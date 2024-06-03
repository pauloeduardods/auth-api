package handlers

import (
	"auth-api/src/internal/domain/auth"
	auth_usecases "auth-api/src/internal/usecases/auth"
	"context"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	useCases *auth_usecases.UseCases
}

func NewAuthHandler(useCases *auth_usecases.UseCases) *AuthHandler {
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
		processRequest(c, loginInput{}, func(ctx context.Context, input loginInput) (*auth.LoginOutput, error) {
			return h.useCases.Login.Execute(ctx, auth_usecases.LoginInput{
				LoginInput: auth.LoginInput{
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
			_, err := h.useCases.ConfirmSignUp.Execute(ctx, auth_usecases.ConfirmSignUpInput{
				ConfirmSignUpInput: auth.ConfirmSignUpInput{
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
		processRequestQuery(c, getMeInput{}, func(ctx context.Context, input getMeInput) (*auth.GetMeOutput, error) {
			return h.useCases.GetMe.Execute(ctx, auth_usecases.GetMeInput{
				GetMeInput: auth.GetMeInput{
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
		processRequest(c, refreshTokenInput{}, func(ctx context.Context, input refreshTokenInput) (*auth.RefreshTokenOutput, error) {
			return h.useCases.RefreshToken.Execute(ctx, auth_usecases.RefreshTokenInput{
				RefreshTokenInput: auth.RefreshTokenInput{
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
		processRequest(c, addMfaInput{}, func(ctx context.Context, input addMfaInput) (*auth.AddMFAOutput, error) {
			return h.useCases.AddMFA.Execute(ctx, auth_usecases.AddMFAInput{
				AddMFAInput: auth.AddMFAInput{
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
		processRequest(c, verifyMfaInput{}, func(ctx context.Context, input verifyMfaInput) (*auth.LoginOutput, error) {
			return h.useCases.VerifyMFA.Execute(ctx, auth_usecases.VerifyMFAInput{
				VerifyMFAInput: auth.VerifyMFAInput{
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
			return h.useCases.AdminRemoveMFA.Execute(ctx, auth_usecases.AdminRemoveMFAInput{
				AdminRemoveMFAInput: auth.AdminRemoveMFAInput{
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
			return h.useCases.RemoveMFA.Execute(ctx, auth_usecases.RemoveMFAInput{
				RemoveMFAInput: auth.RemoveMFAInput{
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
			err := h.useCases.ActivateMFA.Execute(ctx, auth_usecases.ActivateMFAInput{
				ActivateMFAInput: auth.ActivateMFAInput{
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
			err := h.useCases.Logout.Execute(ctx, auth_usecases.LogoutInput{
				LogoutInput: auth.LogoutInput{
					AccessToken: input.AccessToken,
				},
			})
			return err
		})
	}
}

type setPasswordInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Session  string `json:"session"`
}

func (h *AuthHandler) SetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, setPasswordInput{}, func(ctx context.Context, input setPasswordInput) error {
			err := h.useCases.SetPassword.Execute(ctx, auth_usecases.SetPasswordInput{
				SetPasswordInput: auth.SetPasswordInput{
					Username: input.Email,
					Password: input.Password,
					Session:  input.Session,
				},
			})
			return err
		})
	}
}
