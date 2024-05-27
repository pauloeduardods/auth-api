package handlers

import (
	"monitoring-system/server/domain/auth"
	"monitoring-system/server/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth      auth.Auth
	validator validator.Validator
}

func NewAuthHandler(a auth.Auth, validator validator.Validator) *AuthHandler {
	return &AuthHandler{
		auth:      a,
		validator: validator,
	}
}

func (a *AuthHandler) Login() gin.HandlerFunc {
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

func (a *AuthHandler) SignUp() gin.HandlerFunc {
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
