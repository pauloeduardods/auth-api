package handlers

import (
	"monitoring-system/server/domain/auth"
	"monitoring-system/server/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	auth      auth.Auth
	validator validator.Validator
}

type AuthHandler interface {
	Login() gin.HandlerFunc
	SignUp() gin.HandlerFunc
	ConfirmSignUp() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	RefreshToken() gin.HandlerFunc
	CreateAdmin() gin.HandlerFunc
}

func NewAuthHandler(a auth.Auth, validator validator.Validator) AuthHandler {
	return &authHandler{
		auth:      a,
		validator: validator,
	}
}

type loginInput struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

func (a *authHandler) Login() gin.HandlerFunc {
	return func(g *gin.Context) {
		var login loginInput
		if err := g.ShouldBindJSON(&login); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&login)
		if err != nil {
			g.Error(err)
			return
		}

		out, err := a.auth.Login(auth.NewLoginInput(login.Email, login.Password))
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, out)
		}
	}
}

type signUpInput struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
	Name     string `json:"name" binding:"required" validate:"min=3,max=50"`
}

func (a *authHandler) SignUp() gin.HandlerFunc {
	return func(g *gin.Context) {
		var signUp signUpInput

		if err := g.ShouldBindJSON(&signUp); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&signUp)
		if err != nil {
			g.Error(err)
			return
		}

		out, err := a.auth.SignUp(auth.NewSignUpInput(signUp.Email, signUp.Password, signUp.Name))
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, out)
		}
	}
}

type confirmSignUpInput struct {
	Email string `json:"email" binding:"required" validate:"email"`
	Code  string `json:"code" binding:"required" validate:"numeric"`
}

func (a *authHandler) ConfirmSignUp() gin.HandlerFunc {
	return func(g *gin.Context) {
		var confirmSignUp confirmSignUpInput
		if err := g.ShouldBindJSON(&confirmSignUp); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		}

		_, err = a.auth.ConfirmSignUp(auth.NewConfirmSignUpInput(confirmSignUp.Email, confirmSignUp.Code))
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusNoContent, gin.H{})
		}
	}
}

type getUserInput struct {
	AccessToken string `form:"accessToken" binding:"required"`
}

func (a *authHandler) GetUser() gin.HandlerFunc {
	return func(g *gin.Context) {
		var getUser getUserInput
		if err := g.ShouldBindQuery(&getUser); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&getUser)
		if err != nil {
			g.Error(err)
			return
		}

		out, err := a.auth.GetUser(auth.NewGetUserInput(getUser.AccessToken))
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, out)
		}
	}
}

type refreshTokenInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (a *authHandler) RefreshToken() gin.HandlerFunc {
	return func(g *gin.Context) {
		var refreshToken refreshTokenInput
		if err := g.ShouldBindJSON(&refreshToken); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&refreshToken)
		if err != nil {
			g.Error(err)
			return
		}

		out, err := a.auth.RefreshToken(auth.NewRefreshTokenInput(refreshToken.RefreshToken))
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, out)
		}
	}
}

type createAdminInput struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
	Name     string `json:"name" binding:"required" validate:"min=3,max=50"`
}

func (a *authHandler) CreateAdmin() gin.HandlerFunc {
	return func(g *gin.Context) {
		var createAdmin createAdminInput

		if err := g.ShouldBindJSON(&createAdmin); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&createAdmin)
		if err != nil {
			g.Error(err)
			return
		}

		out, err := a.auth.CreateAdmin(auth.NewCreateAdminInput(createAdmin.Email, createAdmin.Password, createAdmin.Name))
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, out)
		}
	}
}
