package handlers

import (
	userServices "monitoring-system/server/internal/domain/user/service"
	"monitoring-system/server/pkg/validator"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	ConfirmSignUp() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	Login() gin.HandlerFunc
	Register() gin.HandlerFunc
}

type AuthHandlerImpl struct {
	authService userServices.AuthService
	validator   validator.Validator
}

func NewAuthHandler(s userServices.AuthService, validator validator.Validator) AuthHandler {
	return &AuthHandlerImpl{
		authService: s,
		validator:   validator,
	}
}

func (a *AuthHandlerImpl) ConfirmSignUp() gin.HandlerFunc {
	return func(g *gin.Context) {
		var confirmSignUp userServices.ConfirmSignUpInput
		if err := g.ShouldBindJSON(&confirmSignUp); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.authService.ConfirmSignUp(confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res)
		}
	}
}

func (a *AuthHandlerImpl) GetUser() gin.HandlerFunc {
	return func(g *gin.Context) {
		accessToken := g.GetHeader("Authorization")
		getUser := userServices.GetUserInput{
			AccessToken: strings.Split(accessToken, " ")[1],
		}

		err := a.validator.Validate(&getUser)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.authService.GetUser(getUser)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res)
		}
	}
}

func (a *AuthHandlerImpl) Login() gin.HandlerFunc {
	return func(g *gin.Context) {
		var login userServices.LoginInput
		if err := g.ShouldBindJSON(&login); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&login)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.authService.Login(login)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res.AuthenticationResult)
		}
	}
}

func (a *AuthHandlerImpl) Register() gin.HandlerFunc {
	return func(g *gin.Context) {
		var signUp userServices.SignUpInput
		if err := g.ShouldBindJSON(&signUp); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&signUp)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.authService.SignUp(signUp)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res)
		}
	}
}
