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
	SignUp(groupName auth.UserGroup) gin.HandlerFunc
	ConfirmSignUp() gin.HandlerFunc
	GetUser() gin.HandlerFunc
}

func NewAuthHandler(a auth.Auth, validator validator.Validator) AuthHandler {
	return &authHandler{
		auth:      a,
		validator: validator,
	}
}

type LoginInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

func (a *authHandler) Login() gin.HandlerFunc {
	return func(g *gin.Context) {
		var login LoginInput
		if err := g.ShouldBindJSON(&login); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&login)
		if err != nil {
			g.Error(err)
			return
		}

		out, err := a.auth.Login(auth.NewLoginInput(login.Username, login.Password))
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, out)
		}
	}
}

type SignUpInput struct {
	Username  string         `json:"username" binding:"required" validate:"email"`
	Password  string         `json:"password" binding:"required" validate:"min=8"`
	Name      string         `json:"name" binding:"required" validate:"min=3,max=50"`
	GroupName auth.UserGroup `json:"groupName" binding:"required"`
}

func (a *authHandler) SignUp(groupName auth.UserGroup) gin.HandlerFunc {
	return func(g *gin.Context) {
		var signUp SignUpInput
		signUp.GroupName = groupName

		if err := g.ShouldBindJSON(&signUp); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&signUp)
		if err != nil {
			g.Error(err)
			return
		}

		out, err := a.auth.SignUp(auth.NewSignUpInput(signUp.Username, signUp.Password, signUp.Name, signUp.GroupName))
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, out)
		}
	}
}

type ConfirmSignUpInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Code     string `json:"code" binding:"required" validate:"numeric"`
}

func (a *authHandler) ConfirmSignUp() gin.HandlerFunc {
	return func(g *gin.Context) {
		var confirmSignUp ConfirmSignUpInput
		if err := g.ShouldBindJSON(&confirmSignUp); err != nil {
			g.Error(err)
			return
		}

		err := a.validator.Validate(&confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		}

		_, err = a.auth.ConfirmSignUp(auth.NewConfirmSignUpInput(confirmSignUp.Username, confirmSignUp.Code))
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusNoContent, gin.H{})
		}
	}
}

type GetUserInput struct {
	AccessToken string `json:"accessToken" form:"accessToken" binding:"required"`
}

func (a *authHandler) GetUser() gin.HandlerFunc {
	return func(g *gin.Context) {
		var getUser GetUserInput
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
