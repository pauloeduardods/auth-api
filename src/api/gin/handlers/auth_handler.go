package handlers

import (
	"auth-api/src/internal/modules/user-manager/domain/admin"
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/internal/modules/user-manager/domain/user"
	auth_usecases "auth-api/src/internal/modules/user-manager/usecases/auth"
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
				Username: input.Email,
				Code:     input.Code,
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
		processRequest(c, setPasswordInput{}, func(ctx context.Context, input setPasswordInput) (*auth.LoginOutput, error) {
			out, err := h.useCases.SetPassword.Execute(ctx, auth_usecases.SetPasswordInput{
				SetPasswordInput: auth.SetPasswordInput{
					Username: input.Email,
					Password: input.Password,
					Session:  input.Session,
				},
			})
			return out, err
		})
	}
}

type addGroupInput struct {
	Email string         `json:"email"`
	Name  *string        `json:"name"`
	Phone *string        `json:"phone"`
	Group auth.UserGroup `json:"group"`
}

func (h *AuthHandler) AddGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, addGroupInput{}, func(ctx context.Context, input addGroupInput) error {
			var adminName, userName, userPhone string
			if input.Name != nil {
				adminName = *input.Name
				userName = *input.Name
			}
			if input.Phone != nil {
				userPhone = *input.Phone
			}

			err := h.useCases.AddGroup.Execute(ctx, auth_usecases.AddGroupInput{
				AddGroupInput: auth.AddGroupInput{
					Username:  input.Email,
					GroupName: input.Group,
				},
				CreateAdminInput: &admin.CreateAdminInput{
					Email: input.Email,
					Name:  adminName,
				},
				CreateUserInput: &user.CreateUserInput{
					Email: input.Email,
					Name:  userName,
					Phone: &userPhone,
				},
			})
			return err
		})
	}
}

type removeGroupInput struct {
	Email string         `json:"email"`
	Group auth.UserGroup `json:"group"`
}

func (h *AuthHandler) RemoveGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, removeGroupInput{}, func(ctx context.Context, input removeGroupInput) error {
			err := h.useCases.RemoveGroup.Execute(ctx, auth_usecases.RemoveGroupInput{
				RemoveGroupInput: auth.RemoveGroupInput{
					Username:  input.Email,
					GroupName: input.Group,
				},
			})
			return err
		})
	}
}

type changePasswordInput struct {
	AccessToken string `json:"accessToken"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (h *AuthHandler) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, changePasswordInput{}, func(ctx context.Context, input changePasswordInput) error {
			err := h.useCases.ChangePassword.Execute(ctx, auth_usecases.ChangePasswordInput{
				AccessToken: input.AccessToken,
				OldPassword: input.OldPassword,
				NewPassword: input.NewPassword,
			})
			return err
		})
	}
}

type resetPasswordInput struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"newPassword"`
}

func (h *AuthHandler) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, resetPasswordInput{}, func(ctx context.Context, input resetPasswordInput) error {
			err := h.useCases.ResetPassword.Execute(ctx, auth_usecases.ResetPasswordInput{
				Username:    input.Email,
				Code:        input.Code,
				NewPassword: input.NewPassword,
			})
			return err
		})
	}
}

type sendForgotPasswordCodeInput struct {
	Email string `json:"email"`
}

func (h *AuthHandler) SendForgotPasswordCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, sendForgotPasswordCodeInput{}, func(ctx context.Context, input sendForgotPasswordCodeInput) error {
			err := h.useCases.SendForgotPasswordCode.Execute(ctx, auth_usecases.SendForgotPasswordCodeInput{
				Username: input.Email,
			})
			return err
		})
	}
}

type sendConfirmationCodeInput struct {
	Email string `json:"email"`
}

func (h *AuthHandler) SendConfirmationCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		processRequestNoOutput(c, sendConfirmationCodeInput{}, func(ctx context.Context, input sendConfirmationCodeInput) error {
			err := h.useCases.SendConfirmationCode.Execute(ctx, auth_usecases.SendConfirmationCodeInput{
				Username: input.Email,
			})
			return err
		})
	}
}
