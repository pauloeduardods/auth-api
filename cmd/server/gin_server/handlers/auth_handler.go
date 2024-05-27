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
}

func NewAuthHandler(a auth.Auth, validator validator.Validator) AuthHandler {
	return &authHandler{
		auth:      a,
		validator: validator,
	}
}

func (a *authHandler) Login() gin.HandlerFunc {
	return func(g *gin.Context) {
		var login auth.LoginInput
		if err := g.ShouldBindJSON(&login); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&login)
		if err != nil {
			g.Error(err)
			return
		}

		out, err := a.auth.Login(login)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, out)
		}
	}
}

func (a *authHandler) SignUp() gin.HandlerFunc {
	return func(g *gin.Context) {
		var signUp auth.SignUpInput
		if err := g.ShouldBindJSON(&signUp); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&signUp)
		if err != nil {
			g.Error(err)
			return
		}

		out, err := a.auth.SignUp(signUp)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, out)
		}
	}
}

func (a *authHandler) ConfirmSignUp() gin.HandlerFunc {
	return func(g *gin.Context) {
		var confirmSignUp auth.ConfirmSignUpInput
		if err := g.ShouldBindJSON(&confirmSignUp); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		}

		_, err = a.auth.ConfirmSignUp(confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusNoContent, gin.H{})
		}
	}
}

func (a *authHandler) GetUser() gin.HandlerFunc {
	return func(g *gin.Context) {
		var getUser auth.GetUserInput
		if err := g.ShouldBindQuery(&getUser); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&getUser)
		if err != nil {
			g.Error(err)
			return
		}

		out, err := a.auth.GetUser(getUser)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, out)
		}
	}
}
